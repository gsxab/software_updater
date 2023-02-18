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

package base

import (
	"context"
	"encoding/json"
	"software_updater/core/action"
	"sync"
)

type Default struct {
}

func (d *Default) Icon() string {
	return "ray-vertex"
}

func (d *Default) InElmNum() int {
	return action.Any
}

func (d *Default) InStrNum() int {
	return action.Any
}

func (d *Default) OutElmNum() int {
	return action.Same
}

func (d *Default) OutStrNum() int {
	return action.Same
}

func (d *Default) Init(context.Context, *sync.WaitGroup) error {
	return nil
}

type DefaultFactory[T any, PT interface {
	action.Action
	*T
}] struct{}

func (r *DefaultFactory[T, PT]) NewAction(args string) (action.Action, error) {
	ret := PT(new(T))
	err := json.Unmarshal([]byte(args), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *DefaultFactory[T, PT]) ToProtoDTO() *action.ProtoDTO {
	t := PT(new(T))
	return t.ToDTO().ProtoDTO
}
