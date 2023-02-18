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
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
)

type RegexpFilter struct {
	base.DefaultFactory[RegexpFilter, *RegexpFilter]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (r *RegexpFilter) Path() action.Path {
	return action.Path{"selector", "filter", "regexp_filter"}
}

func (a *RegexpFilter) Icon() string {
	return "mdi:mdi-filter"
}

func (a *RegexpFilter) InElmNum() int {
	return action.Any
}

func (a *RegexpFilter) InStrNum() int {
	return action.Any
}

func (a *RegexpFilter) OutElmNum() int {
	return 1
}

func (a *RegexpFilter) OutStrNum() int {
	return action.Same
}

func (a *RegexpFilter) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpFilter) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements := input.Elements
	var text string
	for _, element := range elements {
		text, err = element.Text()
		if err != nil {
			logs.Error(ctx, "selenium element get_text failed", err)
			return
		}
		matched := util.Match(a.matcher, a.FullMatch, text)
		if matched {
			output = action.ElementToArgs(element, input)
			return
		}
	}
	err = fmt.Errorf("find matching element failed, matcher: %s, elements: %v", a.Pattern, elements)
	logs.Error(ctx, "element matching failed", err)
	return
}

func (a *RegexpFilter) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"nodes..."},
			Output: []string{"node"},
		},
		Values: map[string]string{"pattern": a.Pattern},
	}
}
