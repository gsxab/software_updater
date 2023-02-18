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

package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"software_updater/core/db/po"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./core/db/dao",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	db, err := gorm.Open(sqlite.Open("./software.db"))
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	g.ApplyBasic(po.Homepage{}, po.Version{}, po.CurrentVersion{})
	g.Execute()

	if db.Migrator().HasTable(&po.Homepage{}) {
		err = db.Migrator().AutoMigrate(po.Homepage{})
	} else {
		err = db.Migrator().CreateTable(&po.Homepage{})
	}
	if err != nil {
		panic(err)
	}

	if db.Migrator().HasTable(&po.Version{}) {
		err = db.Migrator().AutoMigrate(&po.Version{})
	} else {
		err = db.Migrator().CreateTable(&po.Version{})
	}
	if err != nil {
		panic(err)
	}

	if db.Migrator().HasTable(&po.CurrentVersion{}) {
		err = db.Migrator().AutoMigrate(&po.CurrentVersion{})
	} else {
		err = db.Migrator().CreateTable(&po.CurrentVersion{})
	}
	if err != nil {
		panic(err)
	}
}
