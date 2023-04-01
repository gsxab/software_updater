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
	"context"
	"software_updater/core/action"
	"software_updater/core/action/std"
	"software_updater/core/config"
	"software_updater/core/hook"
	"software_updater/core/util/error_util"
)

func DefaultPlugins(config *config.EngineConfig) []Plugin {
	var debugCheckPlugin *HookPlugin
	if config.DebugCheck {
		debugCheckPlugin = getDebugCheckPlugin()
	}

	stdPlugins := []Plugin{
		// basic function
		&ActionPlugin{
			Factories: []action.Factory{
				&std.ReturnEmpty{},
				&std.ReturnConst{},
				&std.AppendConst{},
				&std.Wait{},
				//&action.WaitUntilNext{},
				// browser control
				&std.AccessConst{},
				&std.Access{},
				//&action.MouseMove{},
				&std.Click{},
				// node selector
				&std.CSSSelect{},
				&std.CSSSelectMultiple{},
				&std.CSSSelectAppend{},
				//&action.XPathSelect{},
				//&action.XPathSelectMultiple{},
				//&action.XPathSelectAppend{},
				&std.RegexpFilter{},
				&std.RegexpFilterExtract{},
				//&action.ContainsFilter{},
				// node reader
				&std.ReadText{},
				&std.ReadAttr{},
				// value checker
				//&action.CheckVersion{},
				//&action.CheckDate{},
				// string mutator
				&std.RegexpExtract{},
				&std.Format{},
				&std.ReduceFormat{},
				&std.ReduceJoin{},
				&std.AppendFormat{},
				// url escaper
				&std.URLUnescape{},
				&std.URLEscape{},
				// json handler
				&std.GoJQ{},
				//&std.JsonPath{},
				// encoding
				//&action.Encoding{},
				// base encoder (RFC3548/4648)
				&std.Base64URLDecode{},
				&std.Base64URLEncode{},
				//&action.Base64Decode{},
				//&action.Base64Encode{},
				//&action.Base32Decode{},
				//&action.Base32Encode{},
				//&action.Base32HexDecode{},
				//&action.Base32HexEncode{},
				&std.HexDecode{},
				&std.HexEncode{},
				// compare versions
				&std.CheckVersion{},
				&std.CheckLaterVersion{},
				// mark update
				&std.MarkUpdate{},
				// curl a url
				&std.CURL{},
				&std.CURLSave{},
				// save string
				//&std.SaveString{},
				//&std.SavePage{},
				// store infos
				&std.StoreURL{},
				&std.StoreVersion{},
				&std.StoreDate{},
				&std.StoreDigest{},
				&std.StoreScreenshot{},
				//&action.StoreStr{},
			},
		},
	}
	if debugCheckPlugin != nil {
		stdPlugins = append(stdPlugins, debugCheckPlugin)
	}
	return stdPlugins
}

func getDebugCheckPlugin() *HookPlugin {
	return &HookPlugin{
		RegisterItems: []*hook.RegisterInfo{
			{
				Action: action.Path{action.All},
				Hook: hook.Hook{
					Name: "debug_check_before_node",
					F: func(ctx context.Context, vars *hook.Variables, id string, errs error_util.Collector) {
						a := *vars.ActionPtr
						errs.Collect(action.DynamicCheckInput(a.InElmNum(), len(vars.Input.Elements), id))
						errs.Collect(action.DynamicCheckInput(a.InStrNum(), len(vars.Input.Strings), id))
					},
				},
				Position: hook.First(),
				Event:    "before",
			},
			{
				Action: action.Path{action.All},
				Hook: hook.Hook{
					Name: "debug_check_after_node",
					F: func(ctx context.Context, vars *hook.Variables, id string, errs error_util.Collector) {
						a := *vars.ActionPtr
						errs.Collect(action.DynamicCheckOutput(a.OutElmNum(), len(vars.Input.Elements), len(vars.Output.Elements), id))
						errs.Collect(action.DynamicCheckOutput(a.OutStrNum(), len(vars.Input.Strings), len(vars.Output.Strings), id))
					},
				},
				Position: hook.First(),
				Event:    "after",
			},
		},
	}
}
