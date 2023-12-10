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
	"fmt"
	"software_updater/core/action"
	"software_updater/core/hook"
	"software_updater/core/util"

	"github.com/gsxab/go-logs"
)

type ActionManager struct {
	categories ActionTrie
}

func NewActionManager() *ActionManager {
	actionManager := &ActionManager{}
	actionManager.categories = NewActionTrie()
	return actionManager
}

func (m *ActionManager) Register(factory action.Factory) bool {
	path := factory.Path()

	_, err := m.categories.PutFactLeaf(path, factory)
	if err != nil {
		return false
	}

	return true
}

func (m *ActionManager) RegisterHook(info *hook.RegisterInfo) error {
	position := info.Position
	if position == nil {
		position = &hook.Position{Cmd: hook.LastCmd}
	}
	err := m.categories.PutHook(info.Action, info.Event, info.Hook, position)
	return err
}

func (m *ActionManager) Action(ctx context.Context, storedAction *StoredAction) (action.Action, []*hook.ActionHooks, error) {
	path := action.Path(storedAction.Path)
	args := storedAction.JSON
	if storedAction.Path == nil {
		path = m.categories.GetPath(storedAction.Name)
	}
	if path == nil {
		logs.ErrorM(ctx, "action path is nil", "stored_action", util.ToJSON(storedAction))
		return nil, nil, fmt.Errorf("action path is nil, storedAction: %v", storedAction)
	}
	factory, hooks, err := m.categories.SearchLeafAllHooks(path)
	if err != nil || factory == nil {
		logs.Error(ctx, "tree leaf search failed", err, "path", path, "err", err)
		return nil, nil, fmt.Errorf("action not found, path: %s, error: %w", path, err)
	}
	if len(args) == 0 {
		args = "{}"
	}
	a, err := factory.NewAction(args)
	if err != nil {
		return nil, nil, fmt.Errorf("action creation failed, path: %s, error: %w", path, err)
	}
	return a, hooks, err
}

type StoredAction struct {
	Path []string `json:"path,omitempty"`
	Name string   `json:"name,omitempty"`
	JSON string   `json:"json,omitempty"`
}

type StoredBranch struct {
	Actions []StoredAction `json:"actions,omitempty"`
	Next    []StoredBranch `json:"next,omitempty"`
}

type StoredFlow struct {
	Root StoredBranch `json:"root,omitempty"`
}
