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

package db

import (
	"fmt"
	"log"
	"software_updater/core/config"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(conf *config.DatabaseConfig) (err error) {
	log.Printf("initializing database: %v", conf)

	switch conf.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.DSN), &gorm.Config{})
	default:
		err = fmt.Errorf("database driver not supported: %s", conf.Driver)
	}
	if err != nil {
		return
	}

	migrator := db.Migrator()
	if !migrator.HasTable(&po.Homepage{}) {
		err = migrator.CreateTable(&po.Homepage{})
		if err != nil {
			return
		}
	}
	if !migrator.HasTable(&po.Version{}) {
		err = migrator.CreateTable(&po.Version{})
		if err != nil {
			return
		}
	}
	if !migrator.HasTable(&po.CurrentVersion{}) {
		err = migrator.CreateTable(&po.CurrentVersion{})
		if err != nil {
			return
		}
	}

	dao.SetDefault(db)

	return nil
}

func DB() *gorm.DB {
	return db
}
