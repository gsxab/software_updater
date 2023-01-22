package engine

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/job"
	"software_updater/core/tools/web"
	"software_updater/core/util/error_util"
	"time"
)

type DefaultEngine struct {
	actionManager *ActionManager
	jobManager    *JobManager
	config        *config.EngineConfig
	driver        selenium.WebDriver
}

func (e *DefaultEngine) InitEngine(engineConfig *config.EngineConfig) {
	e.actionManager = NewActionManager()
	e.config = engineConfig
	e.driver = web.Driver()
}

func (e *DefaultEngine) RegisterAction(factory action.Factory) error {
	e.actionManager.Register(factory)
	return nil
}

func (e *DefaultEngine) RegisterHook(registerItem *hook.RegisterInfo) error {
	return e.actionManager.RegisterHook(registerItem)
}

func (e *DefaultEngine) Load(ctx context.Context, homepage *po.Homepage) (*job.Flow, error) {
	flow, err := e.jobManager.Resolve(ctx, homepage.Actions, e.actionManager, e.config)
	return flow, err
}

func (e *DefaultEngine) Crawl(ctx context.Context, hp *po.Homepage) error {
	flow, err := e.Load(ctx, hp)
	if err != nil {
		return err
	}

	version, err := e.jobManager.RunJobs(ctx, flow, e.driver, hp.Current)
	if err != nil {
		return err
	}

	err = e.updateCurrentVersion(ctx, version, hp.Current)
	if err != nil {
		return err
	}

	return nil
}

func (e *DefaultEngine) NeedCrawl(ctx context.Context) (hps []*po.Homepage, err error) {
	hpDAO := dao.Homepage
	cvDAO := dao.CurrentVersion
	if e.config.ForceCrawl {
		hps, err = hpDAO.WithContext(ctx).Find()
	} else {
		hps, err = hpDAO.WithContext(ctx).LeftJoin(cvDAO).
			Where(cvDAO.ScheduledAt.IsNull()).Or(cvDAO.ScheduledAt.Lte(time.Now())).
			Find()
	}
	return
}

func (e *DefaultEngine) CrawlAll(ctx context.Context) error {
	hps, err := e.NeedCrawl(ctx)
	if err != nil {
		return err
	}

	errs := error_util.NewCollector()
	for _, hp := range hps {
		errs.Collect(e.Crawl(ctx, hp))
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
	err = vDAO.WithContext(ctx).Create(v)
	if err != nil {
		return err
	}
	info, err := cvDAO.WithContext(ctx).Where(vDAO.ID.Eq(cv.ID)).UpdateSimple(cvDAO.CurrentVersionID.Value(v.ID))
	if err != nil {
		return err
	}
	if info.Error != nil {
		return info.Error
	}
	if info.RowsAffected != 1 {
		return fmt.Errorf("rows affected unexpected, expected: 1, real: %d", info.RowsAffected)
	}
	return
}
