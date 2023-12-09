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

package po

import (
	"time"

	"gorm.io/gorm"
)

type Version struct {
	gorm.Model
	Name       string          `gorm:"column:name;index:idx_versions_name_version"`
	Version    string          `gorm:"column:version;index:idx_versions_name_version"`
	Filename   *string         `gorm:"column:filename"`
	Picture    *string         `gorm:"column:picture"`
	Link       *string         `gorm:"column:link"`
	Digest     *string         `gorm:"column:digest"`
	RemoteDate *time.Time      `gorm:"column:remote_date;type:date"`
	LocalTime  *time.Time      `gorm:"column:local_time"`
	Previous   *uint           `gorm:"column:previous_version_id;index:idx_versions_previous_version_id"`
	CV         *CurrentVersion `gorm:"foreignKey:VersionID;references:ID"`
}
