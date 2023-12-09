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
	"software_updater/ui/dto"
	"time"

	"github.com/gsxab/logs"
)

func GetList(ctx context.Context) ([]*dto.ListItemDTO, error) {
	hpDAO := dao.Homepage

	hps, err := hpDAO.WithContext(ctx).Preload(hpDAO.Current).Preload(hpDAO.Current.Version).
		Order(hpDAO.Name).Find()
	if err != nil {
		logs.Error(ctx, "list query failed", err)
		return nil, err
	}

	data := make([]*dto.ListItemDTO, 0, len(hps))
	for _, hp := range hps {
		datum := &dto.ListItemDTO{
			Name:    hp.Name,
			PageURL: hp.HomepageURL,
		}
		if hp.Current != nil {
			datum.ScheduledDate = dto.ToDateDTO(hp.Current.ScheduledAt, time.Local)
			if hp.Current.Version != nil {
				datum.Version = &hp.Current.Version.Version
				datum.UpdateDate = dto.ToDateDTO(hp.Current.Version.LocalTime, time.Local)
			}
		}
		data = append(data, datum)
	}

	return data, nil
}
