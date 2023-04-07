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

package flow

import "software_updater/core/action"

type DebugInfo struct {
	Err    error
	Input  *action.Args
	Output *action.Args
}

type StepDTO struct {
	ActionDTO *action.DTO `json:"action"`
	StepName  string      `json:"step_name"`
	State     int         `json:"state"`
	StateDesc string      `json:"state_desc"`
	Duration  *string     `json:"duration,omitempty"`
	DebugInfo *DebugInfo  `json:"debug_info,omitempty"`
}

type BranchDTO struct {
	Steps []*StepDTO   `json:"steps"`
	Next  []*BranchDTO `json:"next,omitempty"`
}

type DTO struct {
	*BranchDTO
	Name    string `json:"name"`
	Version string `json:"version"`
	Desc    string `json:"desc"`
}

type TaskDTO struct {
	*DTO
	State State `json:"state"`
}

func ToStepDTO(step Step) *StepDTO {
	dto := step.ToDTO()
	if dto.ActionDTO == nil {
		dto.ActionDTO = action.ToDTO(step.Action())
	}
	return dto
}

func (f *Flow) ToDTO() *DTO {
	return &DTO{
		BranchDTO: f.makeBranchDTO(f.Root),
		Name:      f.Name,
		Version:   f.Version,
		Desc:      f.Desc,
	}
}

func (f *Flow) makeBranchDTO(b *Branch) *BranchDTO {
	result := &BranchDTO{}

	stepDTOs := make([]*StepDTO, 0, len(b.Steps))
	for _, step := range b.Steps {
		stepDTOs = append(stepDTOs, ToStepDTO(step))
	}
	result.Steps = stepDTOs

	nextDTOs := make([]*BranchDTO, 0, len(b.Next))
	for _, branch := range b.Next {
		nextDTOs = append(nextDTOs, f.makeBranchDTO(branch))
	}
	result.Next = nextDTOs

	return result
}

func ToFlowDTO(flow *Flow) *DTO {
	return flow.ToDTO()
}
