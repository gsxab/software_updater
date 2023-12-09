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
	"fmt"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"

	"github.com/gsxab/logs"
	"github.com/tebeka/selenium"
)

type CSSSelectChild struct {
	base.Default
	base.DefaultFactory[CSSSelectChild, *CSSSelectChild]
	Selector string `json:"selector"`
	Index    int    `json:"index"`
}

func (a *CSSSelectChild) Path() action.Path {
	return action.Path{"browser", "selector", "css", "css_select_child"}
}

func (a *CSSSelectChild) Icon() string {
	return "mdi:mdi-select"
}

func (a *CSSSelectChild) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	if len(input.Elements) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Elements), a.Index)
		logs.Error(ctx, "element slice indexing failed", err, "elements", input.Elements, "index", a.Index)
		return
	}
	parent := input.Elements[a.Index]
	element, err := parent.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = action.IndexedElementToArgs(element, a.Index, input)
	return
}

func (a *CSSSelectChild) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			ReadPage: true,
			Output:   []string{"node"},
		},
		Values: map[string]string{"selector": a.Selector},
	}
}
