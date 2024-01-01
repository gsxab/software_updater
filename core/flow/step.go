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

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"sync"

	"github.com/gsxab/go-error_util/errcollect"
	"github.com/tebeka/selenium"
)

type State int

const (
	// Init means a step or a task is being initialized.
	Init State = iota + 1
	// Pending means a task is pending for the runner.
	Pending
	// Ready means a step is successfully initialized, and never run by the runner.
	Ready
	// Processing means a step or a task is in process.
	Processing
	// Success means a step or a task has exited with success.
	Success
	// Failure means a step fails to be initialized, or a task has exited with failure.
	Failure
	// Cancelled means a step has been cancelled when running, or a task has been cancelled before exiting.
	Cancelled
	// Skipped means a step has marked itself skipped (usually a checker finds nothing to check).
	Skipped
	// Aborted means a step has at least one error reported in running.
	Aborted
	// EarlySuccess means a step has marked itself an early success, or a task has exited with a step marked early success.
	EarlySuccess
)

func (s State) Int() int {
	return int(s)
}

type Step interface {
	SetAction(action action.Action, hooks []*hook.ActionHooks)
	Action() action.Action
	InitAction(ctx context.Context, errs errcollect.Collector, wg *sync.WaitGroup)
	RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs errcollect.Collector, wg *sync.WaitGroup) (output *action.Args, finishBranch bool, earlySuccess bool, stopFlow bool, err error)
	State() State
	SetState(State)
	SetStateDesc(string)
	Name() string
	SetName(string)
	ToDTO() *StepDTO
}
