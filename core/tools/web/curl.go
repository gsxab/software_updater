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
	"io"
	"net/url"
	"os/exec"
	"software_updater/core/config"

	"github.com/tebeka/selenium"
)

func SeleniumCookiesToHeader(url *url.URL, cookies []selenium.Cookie) (result []string) {
	for _, selCookie := range cookies {
		if selCookie.Domain != url.Hostname() {
			continue
		}
		result = append(result, "-H", selCookie.Name+"="+selCookie.Value)
	}
	return
}

func CURL(ctx context.Context, url *url.URL, cookies []selenium.Cookie, output io.Writer) error {
	curlCookies := SeleniumCookiesToHeader(url, cookies)
	args := []string{"-s", "-L"}
	args = append(args, config.Current().CURL.ExtraArgs...)
	args = append(args, curlCookies...)
	args = append(args, url.String())
	cmd := exec.CommandContext(ctx, "curl", args...)
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()
	err = cmd.Start()
	if err != nil {
		return err
	}
	_, err = io.Copy(output, stream)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
