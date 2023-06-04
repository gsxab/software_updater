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
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"sync"
	"time"
)

type TakeAndStoreScreenshot struct {
	base.Default
	base.DefaultFactory[TakeAndStoreScreenshot, *TakeAndStoreScreenshot]
}

func (a *TakeAndStoreScreenshot) Path() action.Path {
	return action.Path{"browser", "reader", "screenshot"}
}

func (a *TakeAndStoreScreenshot) Icon() string {
	return "mdi:mdi-image-plus-outline"
}

func (a *TakeAndStoreScreenshot) OutStrNum() int {
	return action.OneMore
}

func (a *TakeAndStoreScreenshot) getFilename(name string) string {
	encodedName := base64.URLEncoding.EncodeToString([]byte(name))
	dateSuffix := time.Now().Format("2006-01-02")
	return encodedName + "@" + dateSuffix + ".png"
}

func (a *TakeAndStoreScreenshot) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	filename, err := web.TakeScreenshot(ctx, driver, version.Name)
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return
	}

	version.Picture = &filename
	output = input
	return
}

func (a *TakeAndStoreScreenshot) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{},
	}
}
