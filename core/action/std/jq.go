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
	"fmt"
	"github.com/itchyny/gojq"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"
)

type GoJQ struct {
	base.StringMutator
	base.DefaultFactory[GoJQ, *GoJQ]
	Format string `json:"format"`
	query  *gojq.Query
}

func (a *GoJQ) Path() action.Path {
	return action.Path{"string", "mutator", "go_jq"}
}

func (a *GoJQ) Icon() string {
	return "mdi:mdi-code-json"
}

func (a *GoJQ) Init(context.Context, *sync.WaitGroup) (err error) {
	a.query, err = gojq.Parse(a.Format)
	return
}

func (a *GoJQ) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		it := a.query.RunWithContext(ctx, text)
		v, ok := it.Next() // Only read the first one
		if !ok {
			return "", nil
		}
		err, ok := v.(error)
		if ok {
			return "", err
		}
		return fmt.Sprintf("%v", v), err
	})
}
