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

type ReadAttr struct {
	base.Default
	base.DefaultFactory[ReadAttr, *ReadAttr]
	Attribute string `json:"attr"`
}

func (a *ReadAttr) Path() action.Path {
	return action.Path{"browser", "reader", "read_attr"}
}

func (a *ReadAttr) Icon() string {
	return "mdi:mdi-text-box-search"
}

func (a *ReadAttr) InElmNum() int {
	return 1
}

func (a *ReadAttr) InStrNum() int {
	return action.Any
}

func (a *ReadAttr) OutElmNum() int {
	return 1
}

func (a *ReadAttr) OutStrNum() int {
	return 1
}

func (a *ReadAttr) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	text, err := input.Elements[0].GetAttribute(a.Attribute)
	if err != nil {
		logs.Error(ctx, "selenium element get_attr failed", err, "attr", a.Attribute)
		return
	}
	output = action.StringToArgs(text, input)
	return
}

func (a *ReadAttr) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"node"},
			Output: []string{"attribute_text"},
		},
		Values: map[string]string{"attribute": a.Attribute},
	}
}
