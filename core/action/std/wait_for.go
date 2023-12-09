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
	"sync"
	"time"

	"github.com/gsxab/logs"
	"github.com/tebeka/selenium"
)

type Wait struct {
	base.Default
	base.DefaultFactory[Wait, *Wait]
	Delay string `json:"delay"`
	delay time.Duration
}

func (a *Wait) Path() action.Path {
	return action.Path{"basic", "wait", "delay"}
}

func (a *Wait) Icon() string {
	return "mdi:mdi-timer"
}

func (a *Wait) Init(ctx context.Context, _ *sync.WaitGroup) (err error) {
	a.delay, err = time.ParseDuration(a.Delay)
	if err != nil {
		logs.Error(ctx, "duration parsing failed", err, "duration", a.Delay)
		return
	}
	return
}

func (a *Wait) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = input
	t := time.NewTimer(a.delay)
	select {
	case <-ctx.Done():
		logs.WarnM(ctx, "delay action cancelled")
		exit = action.Cancelled
		if !t.Stop() {
			<-t.C
		}
	case <-t.C:
	}
	return
}

func (a *Wait) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"delay": a.Delay},
	}
}
