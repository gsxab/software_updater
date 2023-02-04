package engine

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/action/std_action"
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
				&std_action.ReturnEmpty{},
				&std_action.ReturnConst{},
				&std_action.AppendConst{},
				&std_action.WaitFor{},
				//&action.WaitUntilNext{},
				// browser control
				&std_action.AccessConst{},
				&std_action.Access{},
				//&action.MouseMove{},
				&std_action.Click{},
				// node selector
				&std_action.CSSSelect{},
				&std_action.CSSSelectMultiple{},
				&std_action.CSSSelectAppend{},
				//&action.XPathSelect{},
				//&action.XPathSelectMultiple{},
				//&action.XPathSelectAppend{},
				&std_action.RegexpFilter{},
				//&action.ContainsFilter{},
				// node reader
				&std_action.ReadText{},
				&std_action.ReadAttr{},
				// value checker
				//&action.CheckVersion{},
				//&action.CheckDate{},
				// string mutator
				&std_action.RegexpExtract{},
				&action.Format{},
				&std_action.ReduceFormat{},
				&std_action.AppendFormat{},
				&std_action.URLUnescape{},
				&std_action.URLEscape{},
				// encoding
				//&action.Encoding{},
				// base encoder (RFC3548/4648)
				&std_action.Base64URLDecode{},
				&std_action.Base64URLEncode{},
				//&action.Base64Decode{},
				//&action.Base64Encode{},
				//&action.Base32Decode{},
				//&action.Base32Encode{},
				//&action.Base32HexDecode{},
				//&action.Base32HexEncode{},
				&std_action.HexDecode{},
				&std_action.HexEncode{},
				// curl a url
				&std_action.CURL{},
				&std_action.CURLSave{},
				// store infos
				&std_action.StoreURL{},
				&std_action.StoreVersion{},
				&std_action.StoreDate{},
				&std_action.StoreDigest{},
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
