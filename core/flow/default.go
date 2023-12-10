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
	"fmt"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"sync"
	"time"

	"github.com/gsxab/go-error_util/errcollect"
	"github.com/gsxab/go-logs"
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
)

type DefaultStep struct {
	name      string
	action    action.Action
	hooks     []*hook.ActionHooks
	state     State // Init, Failure; Cancelled; Ready; Processing; Success, Aborted
	start     time.Time
	end       time.Time
	stateDesc string
}

func (j *DefaultStep) SetAction(action action.Action, hooks []*hook.ActionHooks) {
	j.action = action
	j.hooks = hooks
}

func (j *DefaultStep) Action() action.Action {
	return j.action
}

func (j *DefaultStep) InitAction(ctx context.Context, errs errcollect.Collector, wg *sync.WaitGroup) {
	j.state = Init

	select {
	case <-ctx.Done():
		j.state = Cancelled
		return
	default:
	}
	// ready

	// hooks: before init
	for _, actionHooks := range j.hooks {
		hooks := actionHooks.BeforeInit
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
			}, j.name, errs)
		}
	}

	// init action
	err := j.action.Init(ctx, wg)

	// hooks: after init
	for i := len(j.hooks) - 1; i >= 0; i-- {
		hooks := j.hooks[i].AfterInit
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
				ErrorPtr:  &err,
			}, j.name, errs)
		}
	}

	// after init
	if err != nil {
		j.state = Failure
		errs.Collect(fmt.Errorf("action init error: %w", err))
		return
	}
	j.state = Ready
}

func (j *DefaultStep) RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs errcollect.Collector, wg *sync.WaitGroup) (output *action.Args, finishBranch bool, stopFlow bool, err error) {
	j.state = Processing

	select {
	case <-ctx.Done():
		j.state = Cancelled
		return nil, true, true, nil
	default:
	}
	// ready

	j.start = time.Now()
	defer func() {
		j.end = time.Now()
	}()

	// run
	output, result, err := j.run(ctx, driver, args, v, errs, wg)

	// after run
	if err != nil {
		j.state = Aborted
		errs.Collect(fmt.Errorf("action run error: %w", err))
		return output, true, true, err
	}
	switch result {
	case action.Cancelled:
		j.state = Cancelled
		finishBranch = true
	case action.Again:
		return j.RunAction(ctx, nil, args, v, errs, nil)
	case action.StopBranch:
		j.state = Success
		finishBranch = true
	case action.StopFlow:
		j.state = Success
		finishBranch = true
		stopFlow = true
	case action.Finished:
		j.state = Success
	case action.Skipped:
		j.state = Skipped
	default:
		errs.Collect(fmt.Errorf("invalid action state: %d", result))
	}
	return
}

func (j *DefaultStep) run(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version,
	errs errcollect.Collector, wg *sync.WaitGroup) (output *action.Args, result action.Result, err error) {
	defer func() {
		if msg := recover(); msg != nil {
			logs.ErrorM(ctx, "recovered failure", "err", msg)
			output = nil
			err = msg.(error)
		}
	}()

	// hooks: before run
	for _, actionHooks := range j.hooks {
		hooks := actionHooks.BeforeRun
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
				Input:     args,
			}, j.name, errs)
		}
	}

	// run action
	output, result, err = j.action.Do(ctx, driver, args, v, wg)

	// hooks: after run
	for i := len(j.hooks) - 1; i >= 0; i-- {
		hooks := j.hooks[i].AfterRun
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
				Input:     args,
				Output:    output,
				ResultPtr: &result,
				ErrorPtr:  &err,
			}, j.name, errs)
		}
	}
	return
}

func (j *DefaultStep) State() State {
	return j.state
}

func (j *DefaultStep) SetState(state State) {
	j.state = state
}

func (j *DefaultStep) SetStateDesc(s string) {
	j.stateDesc = s
}

func (j *DefaultStep) Name() string {
	return j.name
}

func (j *DefaultStep) SetName(name string) {
	j.name = name
}

func (j *DefaultStep) ToDTO() *StepDTO {
	dto := &StepDTO{
		StepName:  j.name,
		State:     j.state.Int(),
		StateDesc: j.stateDesc,
	}
	if !j.end.IsZero() {
		duration := j.end.Sub(j.start).String()
		dto.Duration = &duration
	}
	return dto
}
