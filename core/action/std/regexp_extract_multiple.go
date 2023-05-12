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
	"github.com/tebeka/selenium"
	"regexp"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type RegexpExtractGroups struct {
	base.Default
	base.DefaultFactory[RegexpExtractGroups, *RegexpExtractGroups]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (a *RegexpExtractGroups) Path() action.Path {
	return action.Path{"string", "mutator", "regexp_extract_groups"}
}

func (a *RegexpExtractGroups) Icon() string {
	return "mdi:mdi-regex"
}

func (a *RegexpExtractGroups) InStrNum() int {
	return 1
}

func (a *RegexpExtractGroups) OutStrNum() int {
	return action.Any
}

func (a *RegexpExtractGroups) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpExtractGroups) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	text := input.Strings[0]
	matched, results := util.MatchExtractMultiple(a.matcher, a.FullMatch, text)
	if !matched {
		err = fmt.Errorf("matching failed, pattern: %s, text: %s", a.Pattern, text)
		return
	}
	output = action.StringsToArgs(results, input)
	return
}

func (a *RegexpExtractGroups) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"extracted_text..."},
		},
		Values: map[string]string{"pattern": a.Pattern, "full_match": util.ToJSON(a.FullMatch)},
	}
}
