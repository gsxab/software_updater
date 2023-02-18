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
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/engine"
	"software_updater/core/logs"
)

func StartFlowByName(ctx context.Context, name string) (engine.TaskID, error) {
	hpDAO := dao.Homepage
	hp, err := hpDAO.WithContext(ctx).Preload(hpDAO.Current).Preload(hpDAO.Current.Version).Where(hpDAO.Name.Eq(name)).Take()
	if err != nil {
		logs.Error(ctx, "homepage query failed", err, "name", name)
		return 0, err
	}
	data, err := StartFlow(ctx, hp)
	return data, err
}

func StartAllFlows(ctx context.Context) error {
	err := engine.Instance().RunAll(ctx)
	return err
}

func StartFlow(ctx context.Context, hp *po.Homepage) (engine.TaskID, error) {
	id, err := engine.Instance().Run(ctx, hp)
	if err != nil {
		return 0, err
	}
	return id, nil
}
