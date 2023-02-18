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
	"software_updater/core/action"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/core/util/error_util"
)

type StringMutator struct {
	Default
	Skip []int `json:"skip,omitempty"`
}

func (a *StringMutator) Icon() string {
	return "mdi:mdi-text-box-edit"
}

func (a *StringMutator) Mutate(input *action.Args, mutate func(text string) string) (output *action.Args, exit action.Result, err error) {
	results := make([]string, 0, len(input.Strings))
	skipIndex := 0
	for index, text := range input.Strings {
		if skipIndex < len(a.Skip) && index == a.Skip[skipIndex] {
			skipIndex++
			results = append(results, text)
			continue
		}

		result := mutate(text)
		results = append(results, result)
	}
	output = action.StringsToArgs(results, input)
	return
}

func (a *StringMutator) MutateWithErr(ctx context.Context, input *action.Args, mutate func(text string) (string, error)) (output *action.Args, exit action.Result, err error) {
	errs := error_util.NewCollector()
	output, exit, err = a.Mutate(input, func(text string) string {
		result, err := mutate(text)
		errs.CollectWithLog(err, func(err error) {
			logs.Error(ctx, "string mutating failed", err, "input", text)
		})
		return result
	})
	errs.Collect(err)
	return output, exit, errs.ToError()
}

func (a *StringMutator) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text..."},
			Output: []string{"formatted_text..."},
		},
		Values: map[string]string{"skip": util.ToJSON(a.Skip)},
	}
}
