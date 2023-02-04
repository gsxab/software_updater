package engine

import (
	"context"
	"software_updater/core/action"
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
				&action.ReturnEmpty{},
				&action.ReturnConst{},
				&action.AppendConst{},
				&action.WaitFor{},
				//&action.WaitUntilNext{},
				// browser control
				&action.AccessConst{},
				&action.Access{},
				//&action.MouseMove{},
				&action.Click{},
				// node selector
				&action.CSSSelect{},
				&action.CSSSelectMultiple{},
				&action.CSSSelectAppend{},
				//&action.XPathSelect{},
				//&action.XPathSelectMultiple{},
				//&action.XPathSelectAppend{},
				&action.RegexpFilter{},
				//&action.ContainsFilter{},
				// node reader
				&action.ReadText{},
				&action.ReadAttr{},
				// value checker
				//&action.CheckVersion{},
				//&action.CheckDate{},
				// string mutator
				&action.RegexpExtract{},
				&action.Format{},
				&action.ReduceFormat{},
				&action.AppendFormat{},
				&action.URLUnescape{},
				&action.URLEscape{},
				// encoding
				//&action.Encoding{},
				// base encoder (RFC3548/4648)
				&action.Base64URLDecode{},
				&action.Base64URLEncode{},
				//&action.Base64Decode{},
				//&action.Base64Encode{},
				//&action.Base32Decode{},
				//&action.Base32Encode{},
				//&action.Base32HexDecode{},
				//&action.Base32HexEncode{},
				&action.HexDecode{},
				&action.HexEncode{},
				// curl a url
				&action.CURL{},
				&action.CURLSave{},
				// store infos
				&action.StoreURL{},
				&action.StoreVersion{},
				&action.StoreDate{},
				&action.StoreDigest{},
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
