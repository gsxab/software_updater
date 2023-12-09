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
	"software_updater/core/util"
	"strings"
	"sync"

	"github.com/tebeka/selenium"
)

type ReduceJoin struct {
	base.StringMutator
	base.DefaultFactory[ReduceJoin, *ReduceJoin]
	Sep string `json:"sep"`
}

func (a *ReduceJoin) Path() action.Path {
	return action.Path{"string", "mutator", "reduce_join"}
}

func (a *ReduceJoin) Icon() string {
	return "mdi:mdi-text-box-plus"
}

func (a *ReduceJoin) OutStrNum() int {
	return 1 + len(a.Skip)
}

func (a *ReduceJoin) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	results := make([]string, 1, len(input.Strings)+1)
	skipIndex := 0
	texts := make([]string, 0, len(input.Strings)+1)
	for index, text := range input.Strings {
		// strings skipped will be pushed into results[1:]
		if skipIndex < len(a.Skip) && index == a.Skip[skipIndex] {
			skipIndex++
			results = append(results, text)
			continue
		}
		// the remaining will be formatted
		texts = append(texts, text)
	}

	result := strings.Join(texts, a.Sep)
	results[0] = result
	output = action.StringsToArgs(results, input)
	return
}

func (a *ReduceJoin) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"formatted_text"},
		},
		Values: map[string]string{"sep": a.Sep, "skip": util.ToJSON(a.Skip)},
	}
}
