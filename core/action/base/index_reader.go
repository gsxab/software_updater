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

package base

import (
	"context"
	"fmt"
	"software_updater/core/action"
	"software_updater/core/logs"
)

type IndexReader struct {
	Default
	Index int `json:"index"`
}

func (a *IndexReader) Read(ctx context.Context, input *action.Args, callback func(text string)) (output *action.Args, exit action.Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	callback(text)
	output = input
	return
}

func (a *IndexReader) ReadDirectly(ctx context.Context, input *action.Args) (text string, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text = input.Strings[a.Index]
	return
}

func (a *IndexReader) ReadWithErr(ctx context.Context, input *action.Args, callback func(text string) error) (output *action.Args, exit action.Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	err = callback(text)
	output = input
	return
}
