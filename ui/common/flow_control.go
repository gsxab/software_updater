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
	"software_updater/core/util"
	"time"

	"github.com/gsxab/go-logs"
)

func StartFlow(ctx context.Context, name string, force bool) (map[string]engine.TaskID, error) {
	hpDAO := dao.Homepage
	query := hpDAO.WithContext(ctx).Preload(hpDAO.Current).Preload(hpDAO.Current.Version)
	if name != "all" {
		query = query.Where(hpDAO.Name.Eq(name))
	} else {
		query = query.Where(hpDAO.NoUpdate.Not())
		if !force {
			cvDAO := dao.CurrentVersion
			query = query.Join(cvDAO, cvDAO.Name.EqCol(hpDAO.Name)).Where(cvDAO.ScheduledAt.Gt(time.Now()))
		}
	}
	hps, err := query.Find()
	if err != nil {
		logs.Error(ctx, "homepage query failed", err, "name", name, "force", force)
		return nil, err
	}
	logs.InfoM(ctx, "homepages found", "name", name, "force", force, "list", util.ToJSON(hps))
	data := make(map[string]engine.TaskID)
	for _, hp := range hps {
		id, err := startFlow(ctx, hp)
		if err != nil {
			return data, err
		}
		data[hp.Name] = id
	}
	return data, nil
}

func startFlow(ctx context.Context, hp *po.Homepage) (engine.TaskID, error) {
	logs.InfoM(ctx, "starting flow", "name", hp.Name)
	id, err := engine.Instance().Run(ctx, hp)
	if err != nil {
		return 0, err
	}
	return id, nil
}
