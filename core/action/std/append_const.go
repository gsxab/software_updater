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
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type AppendConst struct {
	base.Default
	base.DefaultFactory[AppendConst, *AppendConst]
	Val string `json:"val"`
}

func (a *AppendConst) Path() action.Path {
	return action.Path{"basic", "value_generator", "append_value"}
}

func (a *AppendConst) Icon() string {
	return "mdi:mdi-text-box-plus"
}

func (a *AppendConst) OutStrNum() int {
	return action.OneMore
}

func (a *AppendConst) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = action.AnotherStringToArgs(a.Val, input)
	return
}

func (a *AppendConst) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"value": util.ToJSON(a.Val)},
	}
}
