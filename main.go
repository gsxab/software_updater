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
	"context"
	"log"
	"software_updater/core/config"
	"software_updater/core/db"
	"software_updater/core/engine"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/ui/webui"
)

func main() {
	conf, err := config.Load("./conf.yaml")
	if err != nil {
		log.Panic(err)
	}

	err = db.InitDB(conf.Database)
	if err != nil {
		log.Panic(err)
	}

	err = web.InitSelenium(conf.Selenium)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = web.StopSelenium()
	}()

	e, err := engine.InitEngine(conf.Engine)
	if err != nil {
		log.Panic(err)
	}

	uiMode := conf.Extra["ui_mode"]
	switch uiMode {
	case "", "web":
		logs.InfoM(context.Background(), "web ui selected")
		err = webui.InitAndRun(context.Background(), conf.Extra["web_ui_setting"])
		defer engine.DestroyEngine(context.Background())
		if err != nil {
			log.Panic(err)
		}
	case "off":
		err = e.RunAll(context.Background())
		if err != nil {
			log.Panic(err)
		}
	}
}
