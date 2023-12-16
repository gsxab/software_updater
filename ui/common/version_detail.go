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

	"github.com/gsxab/go-logs"
	"github.com/gsxab/go-optional/optional"
)

func GetVersionDetail(ctx context.Context, name string, optionalPage *string, v string) (*dto.VersionDTO, error) {
	flowEnabled := true // default true, except when it is web and no-update
	page, err := optional.New(optionalPage).OrLazyE(func() (string, error) {
		hpDAO := dao.Homepage
		hp, err := hpDAO.WithContext(ctx).Where(hpDAO.Name.Eq(name)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "name", name, "v", v)
			return "", err
		}
		flowEnabled = !hp.NoUpdate
		return hp.HomepageURL, nil
	})
	if err != nil {
		return nil, err
	}

	vDAO := dao.Version

	query := vDAO.WithContext(ctx).Where(vDAO.Name.Eq(name))
	if v == "latest" {
		query = query.Order(vDAO.LocalTime.Desc()).Limit(1)
	} else {
		query = query.Where(vDAO.Version.Eq(v))
	}
	version, err := query.Take()
	if err != nil {
		logs.Error(ctx, "version query failed", err, "name", name, "v", v)
		return nil, err
	}

	data := &dto.VersionDTO{
		Name:        name,
		HomepageURL: page,
		Version:     v,
		PrevVersion: nil,
		NextVersion: nil,
		RemoteDate:  dto.ToDateDTO(version.RemoteDate, time.UTC),
		UpdateDate:  dto.ToDateDTO(version.LocalTime, time.Local),
		Link:        version.Link,
		Digest:      version.Digest,
		Picture:     version.Picture,
		FlowEnabled: flowEnabled,
	}

	if version.Previous != nil {
		previousVersion, err := vDAO.WithContext(ctx).Where(vDAO.ID.Eq(*version.Previous)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "id", *version.Previous)
			return nil, err
		}

		data.PrevVersion = &previousVersion.Version
	}

	nextVersionSlice, err := vDAO.WithContext(ctx).Where(vDAO.Previous.Eq(version.ID)).Find()
	if err != nil {
		logs.Error(ctx, "version query failed", err, "name", name, "v", v)
		return nil, err
	}

	if len(nextVersionSlice) == 1 {
		data.NextVersion = &nextVersionSlice[0].Version
	} else if len(nextVersionSlice) > 1 {
		logs.ErrorM(ctx, "next version more than one", "id", version.ID)
	}

	return data, nil
}
