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
