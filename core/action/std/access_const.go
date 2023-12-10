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

type AccessConst struct {
	base.Default
	base.DefaultFactory[AccessConst, *AccessConst]
	URL string `json:"url"`
}

func (a *AccessConst) Path() action.Path {
	return action.Path{"browser", "access", "goto_const"}
}

func (a *AccessConst) Icon() string {
	return "mdi:mdi-web"
}

func (a *AccessConst) OutElmNum() int {
	return 0
}

func (a *AccessConst) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	base, err := driver.CurrentURL()
	if err != nil {
		logs.Error(ctx, "selenium get current url failed", err)
		return
	}
	url, err := url_util.RelativeURL(base, a.URL)
	if err != nil {
		logs.Error(ctx, "relative url calculation failed", err)
		return
	}
	err = driver.Get(url.String())
	if err != nil {
		logs.Error(ctx, "selenium url access failed", err, "URL", a.URL)
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

func (a *AccessConst) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			OpenPage: true,
		},
		Values: map[string]string{"url": a.URL},
	}
}
