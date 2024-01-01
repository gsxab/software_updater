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

package base

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/db/dao"

	"github.com/gsxab/go-logs"
	version_util "github.com/gsxab/go-version"
)

type VersionComparer struct {
	VersionFilter
}

func (a *VersionComparer) Compare(ctx context.Context, input *action.Args, prevID *uint,
	callback func(prevV *version_util.Version, newV *version_util.Version) (bool, action.Result),
) (output *action.Args, exit action.Result, err error) {
	return a.Filter(ctx, input, func(newVersion *version_util.Version) (res bool, exit action.Result) {
		if prevID == nil {
			logs.InfoM(ctx, "check version skipping: no previous version")
			exit = action.Skipped
			return
		}

		vDAO := dao.Version
		prev, err := vDAO.WithContext(ctx).Where(vDAO.ID.Eq(*prevID)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "id", *prevID)
			return
		}

		previousVersion, err := version_util.Parse(a.VersionFormat, prev.Version)
		if err != nil {
			logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "previous", prev.Version)
			return
		}
		res, exit = callback(previousVersion, newVersion)
		return
	})
}
