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
	"sync"
)

type ReturnEmpty struct {
	base.Default
}

func (a *ReturnEmpty) Path() action.Path {
	return action.Path{"basic", "value_gen", "return_empty"}
}

func (a *ReturnEmpty) Icon() string {
	return "fa:fas fa-solid fa-empty-set"
}

func (a *ReturnEmpty) Do(context.Context, selenium.WebDriver, *action.Args, *po.Version, *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return
}

func (a *ReturnEmpty) ToDTO() *action.DTO {
	return &action.DTO{}
}

func (a *ReturnEmpty) NewAction(_ string) (action.Action, error) {
	return a, nil
}

func (a *ReturnEmpty) ToProtoDTO() *action.ProtoDTO {
	return nil
}
