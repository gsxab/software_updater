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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
	"path"
	"software_updater/core/config"
	"software_updater/core/logs"
	"software_updater/core/util/slice_util"
	"strings"
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

func TakeScreenshot(ctx context.Context, driver selenium.WebDriver, name string) (string, error) {
	filename := base64.URLEncoding.EncodeToString([]byte(name)) + "@" + time.Now().Format("2006-01-02") + ".png"

	width, err := driver.ExecuteScript("return document.body.scrollWidth;", nil)
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return "", err
	}

	height, err := driver.ExecuteScript("return document.body.scrollHeight;", nil)
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return "", err
	}

	err = driver.ResizeWindow("", int(width.(float64)), int(height.(float64)))
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return "", err
	}

	bytes, err := driver.Screenshot()
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return "", err
	}

	err = os.WriteFile(path.Join(config.Current().Files.ScreenshotDir, filename), bytes, os.FileMode(0o644))
	if err != nil {
		logs.Error(ctx, "write file failed", err)
		return "", err
	}

	return filename, nil
}

func ElementToString(element selenium.WebElement) string {
	tag, _ := element.TagName()
	text, _ := element.Text()
	return fmt.Sprintf("<%s>%s", tag, text)
}

func ElementsToString(elements []selenium.WebElement) string {
	return strings.Join(slice_util.Map(elements, ElementToString), ",")
}

type Elements []selenium.WebElement

func (es Elements) String() string {
	return ElementsToString(es)
}

func (es Elements) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.String())
}
