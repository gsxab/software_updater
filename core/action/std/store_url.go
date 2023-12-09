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
	"software_updater/core/util/url_util"
	"sync"

	"github.com/tebeka/selenium"
)

type StoreURL struct {
	base.Default
	base.DefaultFactory[StoreURL, *StoreURL]
	base.IndexReader
}

func (a *StoreURL) Path() action.Path {
	return action.Path{"basic", "value_store", "store_url"}
}

func (a *StoreURL) Icon() string {
	return "mdi:mdi-text-box-check"
}

func (a *StoreURL) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.ReadWithErr(ctx, input, func(text string) error {
		baseURL, err := driver.CurrentURL()
		if err != nil {
			return err
		}
		url, err := url_util.RelativeURL(baseURL, text)
		if err != nil {
			return err
		}
		result := url.String()
		version.Link = &result
		return nil
	})
}

func (a *StoreURL) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"text"},
		},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
