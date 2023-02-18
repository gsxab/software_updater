package engine

import (
	"context"
	"fmt"
	cache "github.com/golang/groupcache/lru"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/job"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util"
	"software_updater/core/util/error_util"
	"time"
)

type DefaultEngine struct {
	actionManager   *ActionManager
	flowInitializer *FlowInitializer
	taskRunner      *TaskRunner
	scheduler       Scheduler
	activeFlows     *cache.Cache // string:*job.Flow
	config          *config.EngineConfig
	driver          selenium.WebDriver
}

func (e *DefaultEngine) InitEngine(ctx context.Context, engineConfig *config.EngineConfig) error {
	e.activeFlows = cache.New(16)
	e.actionManager = NewActionManager()
	e.flowInitializer = NewFlowInitializer()
	e.taskRunner = NewTaskRunner(ctx, true, engineConfig.RunnerCheck, func(ctx context.Context, cv *po.CurrentVersion, v *po.Version) {
		_ = e.updateCurrentVersion(ctx, v, cv)
	})
	e.scheduler = NewScheduler()
	e.config = engineConfig
	e.driver = web.Driver()
	return nil
}

func (e *DefaultEngine) DestroyEngine(ctx context.Context, engineConfig *config.EngineConfig) {
	e.taskRunner.Stop(ctx)
}

func (e *DefaultEngine) RegisterAction(factory action.Factory) error {
	e.actionManager.Register(factory)
	return nil
}

func (e *DefaultEngine) RegisterHook(registerItem *hook.RegisterInfo) error {
	return e.actionManager.RegisterHook(registerItem)
}

func (e *DefaultEngine) Load(ctx context.Context, homepage *po.Homepage, useCache bool) (*job.Flow, error) {
	if useCache {
		if flow, ok := e.activeFlows.Get(homepage.Name); ok {
			return flow.(*job.Flow), nil
		}
	}
	flow, err := e.flowInitializer.Resolve(ctx, homepage.Actions, e.actionManager, e.config)
	if err != nil {
		return nil, err
	}
	err = e.flowInitializer.InitFlow(ctx, flow)
	if err != nil {
		return nil, err
	}
	e.activeFlows.Add(homepage.Name, flow)
	return flow, nil
}

func (e *DefaultEngine) Run(ctx context.Context, hp *po.Homepage) (TaskID, error) {
	flow, err := e.Load(ctx, hp, true)
	if err != nil {
		return 0, err
	}

	id, err := e.taskRunner.EnqueueJob(ctx, flow, hp.Name, hp.Current, hp.HomepageURL)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *DefaultEngine) CheckState(_ context.Context, id TaskID) (bool, job.State, error) {
	return e.taskRunner.GetTaskState(id)
}

func (e *DefaultEngine) NeedCrawl(ctx context.Context) (hps []*po.Homepage, err error) {
	hpDAO := dao.Homepage
	cvDAO := dao.CurrentVersion
	if e.config.ForceCrawl {
		hps, err = hpDAO.WithContext(ctx).Find()
	} else {
		hps, err = hpDAO.WithContext(ctx).Where(hpDAO.NoUpdate.Is(false)).
			Preload(hpDAO.Current).Preload(hpDAO.Current.Version).LeftJoin(cvDAO, cvDAO.Name.EqCol(hpDAO.Name)).
			Where(cvDAO.ScheduledAt.IsNull()).Or(cvDAO.ScheduledAt.Lte(time.Now())).
			Find()
	}
	return
}

func (e *DefaultEngine) RunAll(ctx context.Context) error {
	hps, err := e.NeedCrawl(ctx)
	if err != nil {
		return err
	}

	errs := error_util.NewCollector()
	for _, hp := range hps {
		_, err := e.Run(ctx, hp)
		errs.Collect(err)
	}
	err = errs.ToError()
	if err != nil {
		return err
	}

	return nil
}

func (e *DefaultEngine) updateCurrentVersion(ctx context.Context, v *po.Version, cv *po.CurrentVersion) (err error) {
	vDAO := dao.Version
	cvDAO := dao.CurrentVersion

	if v.LocalTime == nil {
		t := time.Now()
		v.LocalTime = &t
	}

	err = vDAO.WithContext(ctx).Create(v)
	if err != nil {
		logs.Error(ctx, "insert new version failed", err, "v", util.ToJSON(v))
		return err
	}
	schedule := e.scheduler.ScheduleForUpdate(cv, cv.Version, v)
	info, err := cvDAO.WithContext(ctx).Where(cvDAO.ID.Eq(cv.ID)).UpdateSimple(cvDAO.VersionID.Value(v.ID), cvDAO.ScheduledAt.Value(schedule))
	if err != nil {
		logs.Error(ctx, "update current version failed", err, "cv", util.ToJSON(cv), "v.ID", v.ID)
		return err
	}
	if info.Error != nil {
		logs.Error(ctx, "update current version failed", info.Error, "cv", util.ToJSON(cv), "v.ID", v.ID)
		return info.Error
	}
	if info.RowsAffected != 1 {
		err = fmt.Errorf("rows affected unexpected, expected: 1, real: %d", info.RowsAffected)
		logs.Error(ctx, "update current version failed", err, "cv", util.ToJSON(cv), "v.ID", v.ID)
		return err
	}
	return
}

func (e *DefaultEngine) ActionHierarchy(context.Context) (*action.HierarchyDTO, error) {
	return e.actionManager.categories.DTO()
}
