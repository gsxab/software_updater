package engine

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/job"
)

type Engine interface {
	InitEngine(*config.EngineConfig) error
	RegisterAction(factory action.Factory) error
	RegisterHook(registerItem *hook.RegisterInfo) error
	Run(ctx context.Context, homepage *po.Homepage) (TaskID, error)
	CheckState(ctx context.Context, id TaskID) (bool, job.State, error)
	Load(ctx context.Context, homepage *po.Homepage, useCache bool) (*job.Flow, error)
	RunAll(ctx context.Context) error
	ActionHierarchy(ctx context.Context) (*action.HierarchyDTO, error)
}

var engine Engine

func InitEngine(config *config.EngineConfig, extraPlugins ...Plugin) (Engine, error) {
	engine = &DefaultEngine{}
	err := engine.InitEngine(config)
	if err != nil {
		return nil, err
	}

	plugins := DefaultPlugins(config)
	plugins = append(plugins, extraPlugins...)
	for _, plugin := range plugins {
		if plugin == nil {
			continue
		}
		plugin.apply(engine)
	}

	return engine, nil
}

func Instance() Engine {
	return engine
}
