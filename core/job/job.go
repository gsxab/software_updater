/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Copyright (c) 2023. gsxab.
 *
 * This file is part of Software Update Watcher, a.k.a. Zhixin Robot.
 *
 * Software Update Watcher is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
 *
 * Software Update Watcher is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with Software Update Watcher. If not, see <https://www.gnu.org/licenses/>.
 */

package job

import (
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/util/error_util"
	"sync"
)

type State int

const (
	// Init means a job or a task is being initialized.
	Init State = iota + 1
	// Pending means a task is pending for the runner.
	Pending
	// Ready means a job is successfully initialized, and never run by the runner.
	Ready
	// Processing means a job or a task is in process.
	Processing
	// Success means a job or a task has exited with success.
	Success
	// Failure means a job fails to be initialized, or a task has exited with failure.
	Failure
	// Cancelled means a job has been cancelled when running, or a task has been cancelled before exiting.
	Cancelled
	// Skipped means a job has marked itself skipped (usually a checker finds nothing to check).
	Skipped
	// Aborted means a job has at least one error reported in running.
	Aborted
)

func (s State) Int() int {
	return int(s)
}

type Job interface {
	SetAction(action action.Action, hooks []*hook.ActionHooks)
	Action() action.Action
	InitAction(ctx context.Context, errs error_util.Collector, wg *sync.WaitGroup)
	RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) (output *action.Args, finishBranch bool, stopFlow bool, err error)
	State() State
	SetState(State)
	SetStateDesc(string)
	Name() string
	SetName(string)
	ToDTO() *DTO
}

type Branch struct {
	Jobs []Job
	Next []*Branch // reserved for branching
}

type Flow struct {
	Root *Branch
}

type DebugInfo struct {
	Err    error
	Input  *action.Args
	Output *action.Args
}

type DTO struct {
	ActionDTO *action.DTO `json:"action"`
	JobName   string      `json:"job_name"`
	State     int         `json:"state"`
	StateDesc string      `json:"state_desc"`
	Duration  *string     `json:"duration,omitempty"`
	DebugInfo *DebugInfo  `json:"debug_info,omitempty"`
}

type BranchDTO struct {
	Jobs []*DTO       `json:"jobs"`
	Next []*BranchDTO `json:"next,omitempty"`
}

type FlowDTO struct {
	*BranchDTO
}

type TaskDTO struct {
	*FlowDTO
	State State `json:"state"`
}

func ToDTO(job Job) *DTO {
	dto := job.ToDTO()
	if dto.ActionDTO == nil {
		dto.ActionDTO = action.ToDTO(job.Action())
	}
	return dto
}

func (f *Flow) ToDTO() *FlowDTO {
	return &FlowDTO{BranchDTO: f.makeBranchDTO(f.Root)}
}

func (f *Flow) makeBranchDTO(b *Branch) *BranchDTO {
	result := &BranchDTO{}

	jobDTOs := make([]*DTO, 0, len(b.Jobs))
	for _, job := range b.Jobs {
		jobDTOs = append(jobDTOs, ToDTO(job))
	}
	result.Jobs = jobDTOs

	nextDTOs := make([]*BranchDTO, 0, len(b.Next))
	for _, branch := range b.Next {
		nextDTOs = append(nextDTOs, f.makeBranchDTO(branch))
	}
	result.Next = nextDTOs

	return result
}

func ToFlowDTO(flow *Flow) *FlowDTO {
	return flow.ToDTO()
}
