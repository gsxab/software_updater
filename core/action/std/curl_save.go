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
	"os"
	"path"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util/url_util"
	"sync"
)

type CURLSave struct {
	base.Default
}

func (a *CURLSave) Path() action.Path {
	return action.Path{"curl", "access", "curl_save"}
}

func (a *CURLSave) Icon() string {
	return "mdi:mdi-console-network"
}

func (a *CURLSave) InStrNum() int {
	return 2
}

func (a *CURLSave) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = input
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
	selCookies, err := driver.GetCookies()
	if err != nil {
		logs.Error(ctx, "selenium cookies failed", err)
		return
	}

	pathname := path.Join(config.Current().Files.CURLSaveDir, input.Strings[1])
	file, err := os.Open(pathname)
	if err != nil {
		logs.Error(ctx, "file opening failed", err, "filename", input.Strings[1], "pathname", pathname)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logs.Error(ctx, "close file failed", err)
		}
	}(file)

	err = web.CURL(ctx, url, selCookies, file)
	if err != nil {
		logs.Error(ctx, "cURL failed", err, "URL", url)
		return
	}

	return
}

func (a *CURLSave) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: a.ToProtoDTO(),
	}
}

func (a *CURLSave) NewAction(string) (action.Action, error) {
	return a, nil
}

func (a *CURLSave) ToProtoDTO() *action.ProtoDTO {
	return &action.ProtoDTO{
		OpenPage: true,
		Input:    []string{"rel_url", "file_path"},
	}
}
