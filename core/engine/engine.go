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
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/flow"
	"software_updater/core/hook"
)

type Engine interface {
	InitEngine(context.Context, *config.EngineConfig) error
	DestroyEngine(context.Context)

	RegisterAction(factory action.Factory) error
	RegisterHook(registerItem *hook.RegisterInfo) error
	Run(ctx context.Context, homepage *po.Homepage) (TaskID, error)
	CheckState(ctx context.Context, id TaskID) (bool, flow.State, error)
	GetTaskIDMap(ctx context.Context) (map[string]TaskID, error)
	Load(ctx context.Context, homepage *po.Homepage, useCache bool) (*flow.Flow, error)
	RunAll(ctx context.Context) error
	ActionHierarchy(ctx context.Context) (*action.HierarchyDTO, error)

	//RegisterListOp(registerItem *ListOp) error
	//RegisterVersionOp(registerItem *VersionOp) error
	//GetListOps() ([]*ListOp, error)
	//GetVersionOps() ([]*VersionOp, error)
}

var engine Engine

func InitEngine(config *config.EngineConfig, extraPlugins ...Plugin) (Engine, error) {
	engine = &DefaultEngine{}
	ctx := context.Background()
	err := engine.InitEngine(ctx, config)
	if err != nil {
		return nil, err
	}

	plugins := DefaultPlugins(config)
	plugins = append(plugins, extraPlugins...)
	for _, plugin := range plugins {
		if plugin == nil {
			continue
		}
		plugin.apply(engine)
	}

	return engine, nil
}

func DestroyEngine(ctx context.Context) {
	engine.DestroyEngine(ctx)
}

func Instance() Engine {
	return engine
}
