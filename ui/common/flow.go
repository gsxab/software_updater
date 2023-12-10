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
	"software_updater/core/flow"
	"software_updater/ui/dto"

	"github.com/gsxab/go-logs"
)

func GetFlowByName(ctx context.Context, name string, reload bool) (*dto.FlowDTO, error) {
	hpDAO := dao.Homepage
	hp, err := hpDAO.WithContext(ctx).Where(hpDAO.Name.Eq(name)).Take()
	if err != nil {
		logs.Error(ctx, "homepage query failed", err, "name", name)
		return nil, err
	}
	data, err := GetFlow(ctx, hp, reload)
	return data, err
}

func GetFlow(ctx context.Context, hp *po.Homepage, reload bool) (*dto.FlowDTO, error) {
	fl, err := engine.Instance().Load(ctx, hp, !reload)
	if err != nil {
		return nil, err
	}
	data := flow.ToFlowDTO(fl)
	return data, nil
}
