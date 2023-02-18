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

package hook

import (
	"fmt"
	"software_updater/core/util/slice_util"
)

type Event string

const (
	BeforeInitEvent Event = "before_init"
	AfterInitEvent  Event = "after_init"
	BeforeRunEvent  Event = "before"
	AfterRunEvent   Event = "after"
)

type ActionHooks struct {
	BeforeInit []Hook
	AfterInit  []Hook
	BeforeRun  []Hook
	AfterRun   []Hook
}

func (h *ActionHooks) HookPtr(event Event) *[]Hook {
	switch event {
	case BeforeInitEvent:
		return &h.BeforeInit
	case AfterInitEvent:
		return &h.AfterInit
	case BeforeRunEvent:
		return &h.BeforeRun
	case AfterRunEvent:
		return &h.AfterRun
	default:
		panic("unexpected event")
	}
}

func (h *ActionHooks) Get(event Event) []Hook {
	return *h.HookPtr(event)
}

func (h *ActionHooks) PutAt(event Event, hook Hook, pos *Position) error {
	ptr := h.HookPtr(event)

	switch pos.Cmd {
	case FirstCmd:
		*ptr = slice_util.Prepend(*ptr, hook)
	case LastCmd:
		*ptr = append(*ptr, hook)
	case PrevCmd:
		in, idx := slice_util.LinearSearchWithPtr(*ptr, func(x *Hook) bool { return x.Name == pos.Ref })
		if !in {
			return fmt.Errorf("registered not found")
		}
		*ptr = slice_util.Insert(*ptr, idx, hook)
	case NextCmd:
		in, idx := slice_util.LinearSearchWithPtr(*ptr, func(x *Hook) bool { return x.Name == pos.Ref })
		if !in {
			return fmt.Errorf("registered not found")
		}
		*ptr = slice_util.Insert(*ptr, idx+1, hook)
	default:
		return fmt.Errorf("unknown position command")
	}
	return nil
}
