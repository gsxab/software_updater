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
	"software_updater/core/hook"
	"software_updater/core/job"
	"software_updater/core/util/error_util"
	"strconv"
	"sync"
)

type FlowInitializer struct {
}

func NewFlowInitializer() *FlowInitializer {
	return &FlowInitializer{}
}

func (t *FlowInitializer) Resolve(ctx context.Context, actionStr string, actionManager *ActionManager, config *config.EngineConfig) (*job.Flow, error) {
	var storedFlow StoredFlow
	errs := error_util.NewCollector()
	errs.Collect(json.Unmarshal([]byte(actionStr), &storedFlow))

	flow := &job.Flow{
		Root: t.resolveBranch(ctx, storedFlow.Root, actionManager, config, errs, ""),
	}

	return flow, errs.ToError()
}

func (t *FlowInitializer) resolveBranch(ctx context.Context, storedBranch StoredBranch, actionManager *ActionManager,
	config *config.EngineConfig, errs *error_util.ErrorCollector, prefix string) *job.Branch {
	jobs := make([]job.Job, 0, len(storedBranch.Actions))
	for i, storedAction := range storedBranch.Actions {
		a, hooks, err := actionManager.Action(ctx, &storedAction)
		errs.Collect(err)
		if a == nil {
			continue
		}
		j := t.NewJob(ctx, config, a, hooks)
		j.SetName(a.Path().Name() + prefix + "-" + strconv.Itoa(i))
		jobs = append(jobs, j)
	}

	next := make([]*job.Branch, 0, len(storedBranch.Next))
	for i, child := range storedBranch.Next {
		branch := t.resolveBranch(ctx, child, actionManager, config, errs, prefix+"-b"+strconv.Itoa(i))
		next = append(next, branch)
	}

	return &job.Branch{
		Jobs: jobs,
		Next: next,
	}
}

func (t *FlowInitializer) NewJob(_ context.Context, config *config.EngineConfig, action action.Action, hooks []*hook.ActionHooks) job.Job {
	var item job.Job
	if config.DebugLog {
		item = &job.LoggedJob{}
	} else {
		item = &job.DefaultJob{}
	}
	item.SetAction(action, hooks)
	return item
}

func (t *FlowInitializer) StaticCheckFlow(flow *job.Flow) error {
	var elmArgN, strArgN int
	errs := error_util.NewCollector()
	t.staticCheckBranch(flow.Root, elmArgN, strArgN, errs)
	return errs.ToError()
}

func (t *FlowInitializer) staticCheckBranch(thread *job.Branch, elmArgN int, strArgN int, errs *error_util.ErrorCollector) {
	var err error
	for _, j := range thread.Jobs {
		elmArgN, err = action.StaticCheckInput(j.Action().InElmNum(), j.Action().OutElmNum(), elmArgN, j.Name())
		errs.Collect(err)
		strArgN, err = action.StaticCheckInput(j.Action().InStrNum(), j.Action().OutStrNum(), strArgN, j.Name())
		errs.Collect(err)
	}

	for _, child := range thread.Next {
		t.staticCheckBranch(child, elmArgN, strArgN, errs)
	}
}

func (t *FlowInitializer) InitFlow(ctx context.Context, flow *job.Flow) error {
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

	t.initBranch(ctx, flow.Root, &error_util.ChannelCollector{Channel: errChan}, &wg)

	wg.Wait()
	close(errChan)
	<-overChan
	return errs.ToError()
}

func (t *FlowInitializer) initBranch(ctx context.Context, branch *job.Branch, errs error_util.Collector, wg *sync.WaitGroup) {
	for _, j := range branch.Jobs {
		j.InitAction(ctx, errs, wg)
	}

	for _, child := range branch.Next {
		t.initBranch(ctx, child, errs, wg)
	}
}
