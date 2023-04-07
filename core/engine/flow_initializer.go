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

package engine

import (
	"context"
	"encoding/json"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/flow"
	"software_updater/core/hook"
	"software_updater/core/util/error_util"
	"strconv"
	"sync"
)

type FlowInitializer struct {
}

func NewFlowInitializer() *FlowInitializer {
	return &FlowInitializer{}
}

func (t *FlowInitializer) Resolve(ctx context.Context, actionStr string, actionManager *ActionManager) (*flow.Flow, error) {
	var storedFlow StoredFlow
	errs := error_util.NewCollector()
	errs.Collect(json.Unmarshal([]byte(actionStr), &storedFlow))

	fl := &flow.Flow{
		Root: t.resolveBranch(ctx, storedFlow.Root, actionManager, errs, ""),
	}

	return fl, errs.ToError()
}

func (t *FlowInitializer) resolveBranch(ctx context.Context, storedBranch StoredBranch, actionManager *ActionManager,
	errs *error_util.ErrorCollector, prefix string) *flow.Branch {
	steps := make([]flow.Step, 0, len(storedBranch.Actions))
	for i, storedAction := range storedBranch.Actions {
		a, hooks, err := actionManager.Action(ctx, &storedAction)
		errs.Collect(err)
		if a == nil {
			continue
		}
		step := t.NewStep(ctx, a, hooks)
		step.SetName(a.Path().Name() + prefix + "-" + strconv.Itoa(i))
		steps = append(steps, step)
	}

	next := make([]*flow.Branch, 0, len(storedBranch.Next))
	for i, child := range storedBranch.Next {
		branch := t.resolveBranch(ctx, child, actionManager, errs, prefix+"-b"+strconv.Itoa(i))
		next = append(next, branch)
	}

	return &flow.Branch{
		Steps: steps,
		Next:  next,
	}
}

func (t *FlowInitializer) NewStep(_ context.Context, action action.Action, hooks []*hook.ActionHooks) flow.Step {
	var item flow.Step
	if config.Current().Engine.DebugLog {
		item = &flow.LoggedStep{}
	} else {
		item = &flow.DefaultStep{}
	}
	item.SetAction(action, hooks)
	return item
}

func (t *FlowInitializer) StaticCheckFlow(flow *flow.Flow) error {
	var elmArgN, strArgN int
	errs := error_util.NewCollector()
	t.staticCheckBranch(flow.Root, elmArgN, strArgN, errs)
	return errs.ToError()
}

func (t *FlowInitializer) staticCheckBranch(thread *flow.Branch, elmArgN int, strArgN int, errs *error_util.ErrorCollector) {
	var err error
	for _, j := range thread.Steps {
		elmArgN, err = action.StaticCheckInput(j.Action().InElmNum(), j.Action().OutElmNum(), elmArgN, j.Name())
		errs.Collect(err)
		strArgN, err = action.StaticCheckInput(j.Action().InStrNum(), j.Action().OutStrNum(), strArgN, j.Name())
		errs.Collect(err)
	}

	for _, child := range thread.Next {
		t.staticCheckBranch(child, elmArgN, strArgN, errs)
	}
}

func (t *FlowInitializer) InitFlow(ctx context.Context, fl *flow.Flow) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, 16)
	overChan := make(chan struct{})

	errs := error_util.NewCollector()
	go func() {
		err, ok := <-errChan
		for ok {
			errs.Collect(err)
			err, ok = <-errChan
		}
		close(overChan)
	}()

	t.initBranch(ctx, fl.Root, &error_util.ChannelCollector{Channel: errChan}, &wg)

	wg.Wait()
	close(errChan)
	<-overChan
	return errs.ToError()
}

func (t *FlowInitializer) initBranch(ctx context.Context, branch *flow.Branch, errs error_util.Collector, wg *sync.WaitGroup) {
	for _, step := range branch.Steps {
		step.InitAction(ctx, errs, wg)
	}

	for _, child := range branch.Next {
		t.initBranch(ctx, child, errs, wg)
	}
}
