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
	"regexp"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"

	"github.com/gsxab/go-logs"
	"github.com/tebeka/selenium"
)

type RegexpFilterExtract struct {
	base.DefaultFactory[RegexpFilterExtract, *RegexpFilterExtract]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (r *RegexpFilterExtract) Path() action.Path {
	return action.Path{"selector", "filter", "regexp_filter_extract"}
}

func (a *RegexpFilterExtract) Icon() string {
	return "mdi:mdi-regex"
}

func (a *RegexpFilterExtract) InElmNum() int {
	return action.Any
}

func (a *RegexpFilterExtract) InStrNum() int {
	return action.Any
}

func (a *RegexpFilterExtract) OutElmNum() int {
	return 1
}

func (a *RegexpFilterExtract) OutStrNum() int {
	return 1
}

func (a *RegexpFilterExtract) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpFilterExtract) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements := input.Elements
	var text string
	for _, element := range elements {
		text, err = element.Text()
		if err != nil {
			logs.Error(ctx, "selenium element get_text failed", err)
			return
		}
		matched, results := util.MatchExtractMultiple(a.matcher, a.FullMatch, text)
		if matched {
			output = action.StringsToArgs(results, input)
			output.Elements = []selenium.WebElement{element}
			return
		}
	}
	err = fmt.Errorf("find matching element failed, matcher: %s, elements: %v", a.Pattern, elements)
	logs.Error(ctx, "element matching failed", err)
	return
}

func (a *RegexpFilterExtract) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"nodes..."},
			Output: []string{"text"},
		},
		Values: map[string]string{"pattern": a.Pattern},
	}
}
