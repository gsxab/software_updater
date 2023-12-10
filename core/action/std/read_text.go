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

package std

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"

	"github.com/gsxab/go-logs"
	"github.com/tebeka/selenium"
)

type ReadText struct {
	base.Default
}

func (a *ReadText) Path() action.Path {
	return action.Path{"browser", "reader", "read"}
}

func (a *ReadText) Icon() string {
	return "mdi:mdi-text-box-search"
}

func (a *ReadText) InElmNum() int {
	return 1
}

func (a *ReadText) OutElmNum() int {
	return 1
}

func (a *ReadText) OutStrNum() int {
	return 1
}

func (a *ReadText) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	text, err := input.Elements[0].Text()
	if err != nil {
		logs.Error(ctx, "selenium element get_text failed", err)
		return
	}
	output = action.StringToArgs(text, input)
	return
}

func (a *ReadText) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: a.ToProtoDTO(),
	}
}

func (a *ReadText) NewAction(string) (action.Action, error) {
	return a, nil
}

func (a *ReadText) ToProtoDTO() *action.ProtoDTO {
	return &action.ProtoDTO{
		Input:  []string{"node"},
		Output: []string{"text"},
	}
}
