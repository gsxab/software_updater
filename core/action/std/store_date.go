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
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
	"time"
)

type StoreDate struct {
	base.Default
	base.DefaultFactory[StoreDate, *StoreDate]
	base.IndexReader
	Format string `json:"format"`
}

func (a *StoreDate) Path() action.Path {
	return action.Path{"basic", "value_store", "store_date"}
}

func (a *StoreDate) Icon() string {
	return "mdi:mdi-clock-check"
}

func (a *StoreDate) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.ReadWithErr(ctx, input, func(text string) error {
		t, err := time.Parse(a.Format, text) // implicitly use UTC. only date to be shown, and we don't know the right time zones.
		if err != nil {
			logs.Error(ctx, "date parsing failed", err, "text", text)
			return err
		}
		version.RemoteDate = &t
		return nil
	})
}

func (a *StoreDate) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"text"},
		},
		Values: map[string]string{"index": util.ToJSON(a.Index), "format": a.Format},
	}
}
