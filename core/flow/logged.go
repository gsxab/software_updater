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

package flow

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"sync"

	"github.com/gsxab/go-error_util/errcollect"
	"github.com/tebeka/selenium"
)

type LoggedStep struct {
	DefaultStep
	info *DebugInfo
}

func (j *LoggedStep) RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs errcollect.Collector, wg *sync.WaitGroup) (*action.Args, bool, bool, bool, error) {
	output, stop, earlySuccess, cancel, err := j.DefaultStep.RunAction(ctx, driver, args, v, errs, wg)
	j.info = &DebugInfo{
		Err:    err,
		Input:  args,
		Output: output,
	}
	return output, stop, earlySuccess, cancel, err
}

func (j *LoggedStep) ToDTO() *StepDTO {
	dto := j.DefaultStep.ToDTO()
	dto.DebugInfo = j.info
	return dto
}
