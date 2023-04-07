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
	cache "github.com/gsxab/go-generic_lru"
	cache_impl "github.com/gsxab/go-generic_lru/lru_with_rw_lock"
	"software_updater/core/flow"
)

type FlowManager struct {
	activeFlows     cache.Cache[string, *flow.Flow]
	flowInitializer *FlowInitializer
}

func NewFlowManager() *FlowManager {
	fm := &FlowManager{}
	fm.activeFlows = cache_impl.New[string, *flow.Flow](16)
	fm.flowInitializer = NewFlowInitializer()
	return fm
}

func (m *FlowManager) Load(ctx context.Context, name string, val string, actionManager *ActionManager, useCache bool) (*flow.Flow, error) {
	if useCache {
		if fl, ok := m.activeFlows.Get(name); ok {
			return fl, nil
		}
	}
	fl, err := m.flowInitializer.Resolve(ctx, val, actionManager)
	if err != nil {
		return nil, err
	}
	err = m.flowInitializer.InitFlow(ctx, fl)
	if err != nil {
		return nil, err
	}
	m.activeFlows.Add(fl.Name, fl)
	return fl, nil
}
