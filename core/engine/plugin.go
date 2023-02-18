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

package engine

import (
	"software_updater/core/action"
	"software_updater/core/hook"
)

type Plugin interface {
	apply(Engine)
}

type ActionPlugin struct {
	Factories []action.Factory
}

func (p *ActionPlugin) apply(engine Engine) {
	for _, factory := range p.Factories {
		if factory == nil {
			continue
		}
		_ = engine.RegisterAction(factory)
	}
}

type HookPlugin struct {
	RegisterItems []*hook.RegisterInfo
}

func (p *HookPlugin) apply(engine Engine) {
	for _, registerItem := range p.RegisterItems {
		if registerItem == nil {
			continue
		}
		_ = engine.RegisterHook(registerItem)
	}
}
