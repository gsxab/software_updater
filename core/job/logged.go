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

package job

import (
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/util/error_util"
	"sync"
)

type LoggedJob struct {
	DefaultJob
	info *DebugInfo
}

func (j *LoggedJob) RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) (*action.Args, bool, bool, error) {
	output, stop, cancel, err := j.DefaultJob.RunAction(ctx, driver, args, v, errs, wg)
	j.info = &DebugInfo{
		Err:    err,
		Input:  args,
		Output: output,
	}
	return output, stop, cancel, err
}

func (j *LoggedJob) ToDTO() *DTO {
	dto := j.DefaultJob.ToDTO()
	dto.DebugInfo = j.info
	return dto
}
