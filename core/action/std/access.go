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
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/util/url_util"
	"sync"

	"github.com/gsxab/go-logs"
	"github.com/tebeka/selenium"
)

type Access struct {
	base.Default
	base.DefaultFactory[Access, *Access]
}

func (a *Access) Path() action.Path {
	return action.Path{"browser", "access", "goto_url"}
}

func (a *Access) Icon() string {
	return "mdi:mdi-web"
}

func (a *Access) OutElmNum() int {
	return 0
}

func (a *Access) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	relURL := input.Strings[0]
	baseURL, err := driver.CurrentURL()
	if err != nil {
		logs.Error(ctx, "selenium current url failed", err)
		return
	}
	url, err := url_util.RelativeURL(baseURL, relURL)
	if err != nil {
		logs.Error(ctx, "relative url calculation failed", err, "baseURL", baseURL, "relURL", relURL)
		return
	}
	err = driver.Get(url.String())
	if err != nil {
		logs.Error(ctx, "selenium get url failed", err, "baseURL", baseURL, "relURL", relURL)
		return
	}
	size := config.Current().Selenium.WindowSize
	err = driver.ResizeWindow("", size.Width, size.Height)
	if err != nil {
		logs.Error(ctx, "selenium resize failed", err)
		return
	}
	output = action.ElementsToArgs([]selenium.WebElement{}, input)
	return
}

func (a *Access) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			OpenPage: true,
			Input:    []string{"rel_url"},
		},
	}
}
