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

	version_util "github.com/gsxab/go-version"
	"github.com/tebeka/selenium"
)

type CheckLaterVersion struct {
	base.DefaultFactory[CheckLaterVersion, *CheckLaterVersion]
	base.VersionComparer
}

func (a *CheckLaterVersion) Path() action.Path {
	return action.Path{"basic", "value_check", "version_gt"}
}

func (a *CheckLaterVersion) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Compare(ctx, input, version.Previous, func(prevV *version_util.Version, newV *version_util.Version) (bool, action.Result) {
		if prevV.LT(newV) {
			return true, action.Finished // go on for a new version
		}
		return false, action.EarlySuccessBranch // no new version, mark a early success
	})
}

func (a *CheckLaterVersion) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
