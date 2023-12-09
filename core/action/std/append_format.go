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

	"github.com/tebeka/selenium"
)

type AppendFormat struct {
	base.Default
	base.DefaultFactory[AppendFormat, *AppendFormat]
	Format string `json:"format"`
}

func (a *AppendFormat) Path() action.Path {
	return action.Path{"string", "mutator", "append_format"}
}

func (a *AppendFormat) Icon() string {
	return "mdi:mdi-text-box-plus"
}

func (a *AppendFormat) OutStrNum() int {
	return action.OneMore
}

func (a *AppendFormat) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	texts := make([]any, 0, len(input.Strings))
	for _, text := range input.Strings {
		texts = append(texts, text)
	}
	result := fmt.Sprintf(a.Format, texts...)
	output = action.AnotherStringToArgs(result, input)
	return
}

func (a *AppendFormat) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"formatted_text"},
		},
		Values: map[string]string{"pattern": a.Format},
	}
}
