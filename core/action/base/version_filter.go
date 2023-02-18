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
	"software_updater/core/logs"
	"software_updater/core/util/version_util"
)

type VersionFilter struct {
	VersionFormat string `json:"format"`
	IndexReader
}

func (a *VersionFilter) Icon() string {
	return "mdi:mdi-alpha-v-circle"
}

func (a *VersionFilter) Filter(ctx context.Context, input *action.Args,
	callback func(v *version_util.Version) (bool, action.Result),
) (output *action.Args, exit action.Result, err error) {
	output = input

	versionStr, err := a.ReadDirectly(ctx, input)
	if err != nil {
		return
	}

	currentVersion, err := version_util.Parse(a.VersionFormat, versionStr)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "versionStr", versionStr)
		return
	}
	res, exit := callback(currentVersion)
	if exit == action.Skipped {
		return
	}
	if !res {
		logs.InfoM(ctx, "version checker stopping task", "current", currentVersion)
		exit = action.StopFlow
	}
	return
}
