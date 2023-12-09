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

package std

import (
	"context"
	"encoding/base64"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"

	"github.com/tebeka/selenium"
)

type Base64URLEncode struct {
	base.StringMutator
	base.DefaultFactory[Base64URLEncode, *Base64URLEncode]
}

func (a *Base64URLEncode) Path() action.Path {
	return action.Path{"encoder", "rfc4648", "url_base64_encode"}
}

func (a *Base64URLEncode) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Mutate(input, func(text string) string {
		return base64.URLEncoding.EncodeToString([]byte(text))
	})
}
