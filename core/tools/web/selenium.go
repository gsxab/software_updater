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

package web

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"software_updater/core/config"
	"time"
)

var service *selenium.Service
var driver selenium.WebDriver

func InitSelenium(conf *config.SeleniumConfig) (err error) {
	log.Printf("initializing selenium: %v", conf)

	// Run Chrome browser
	service, err = selenium.NewChromeDriverService(conf.DriverPath, 4444)
	if err != nil {
		return err
	}

	prefs := make(map[string]interface{})
	prefs["download_restrictions"] = 3

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: conf.Params, Prefs: prefs})

	driver, err = selenium.NewRemote(caps, "")
	if err != nil {
		return err
	}

	err = driver.SetImplicitWaitTimeout(time.Second * 60)
	if err != nil {
		return err
	}
	err = driver.SetPageLoadTimeout(time.Second * 60)
	if err != nil {
		return err
	}

	return nil
}

func Driver() selenium.WebDriver {
	return driver
}

func StopSelenium() error {
	err := driver.Quit()
	if err != nil {
		return err
	}
	err = service.Stop()
	return err
}
