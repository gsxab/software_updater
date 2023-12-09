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

type CSSSelectMultiple struct {
	base.Default
	base.DefaultFactory[CSSSelectMultiple, *CSSSelectMultiple]
	Selector string `json:"selector"`
}

func (a *CSSSelectMultiple) Path() action.Path {
	return action.Path{"browser", "selector", "css", "css_select_multiple"}
}

func (a *CSSSelectMultiple) Icon() string {
	return "mdi:mdi-select-multiple"
}

func (a *CSSSelectMultiple) OutElmNum() int {
	return action.Any
}

func (a *CSSSelectMultiple) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements, err := driver.FindElements(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	if len(elements) == 0 {
		err = fmt.Errorf("selenium element not found")
		logs.Error(ctx, "selenium element not found", err, "selector", a.Selector)
		return
	}
	output = action.ElementsToArgs(elements, input)
	return
}

func (a *CSSSelectMultiple) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			ReadPage: true,
			Output:   []string{"nodes..."},
		},
		Values: map[string]string{"selector": a.Selector},
	}
}
