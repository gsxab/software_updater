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

package hook

import (
	"context"
	"software_updater/core/action"

	"github.com/gsxab/go-error_util/errcollect"
)

type Variables struct {
	ActionPtr *action.Action
	Input     *action.Args
	Output    *action.Args
	ResultPtr *action.Result
	ErrorPtr  *error
}

type Hook struct {
	F    func(ctx context.Context, vars *Variables, id string, errs errcollect.Collector)
	Name string
}

type RegisterInfo struct {
	Action   action.Path
	Hook     Hook
	Position *Position
	Event    Event
}
