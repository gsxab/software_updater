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
	InitEngine(*config.EngineConfig)
	RegisterAction(factory action.Factory) error
	RegisterHook(registerItem *hook.RegisterInfo) error
	Crawl(ctx context.Context, homepage *po.Homepage) error
	Load(ctx context.Context, homepage *po.Homepage) (*job.Flow, error)
	CrawlAll(ctx context.Context) error
}

var engine Engine

func InitEngine(config *config.EngineConfig, extraPlugins ...Plugin) (Engine, error) {
	engine = &DefaultEngine{}
	engine.InitEngine(config)

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
