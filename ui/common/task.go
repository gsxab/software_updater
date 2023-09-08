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

package common

import (
	"context"
	"software_updater/core/engine"
	"software_updater/ui/dto"
)

func GetTaskStateByID(ctx context.Context, taskID int64) (bool, int, error) {
	exist, state, err := engine.Instance().CheckState(ctx, taskID)
	return exist, int(state), err
}

func GetTaskIDMap(ctx context.Context) (map[string]int64, error) {
	idMap, err := engine.Instance().GetTaskIDMap(ctx)
	return idMap, err
}

func GetTaskMeta(ctx context.Context, taskID int64) (bool, *dto.TaskMetaDTO, error) {
	exists, taskMeta, err := engine.Instance().GetTaskMeta(ctx, taskID)
	return exists, taskMeta, err
}

func GetTaskMetaList(ctx context.Context) ([]*dto.TaskMetaDTO, error) {
	taskMetaList, err := engine.Instance().GetTaskMetaList(ctx)
	return taskMetaList, err
}
