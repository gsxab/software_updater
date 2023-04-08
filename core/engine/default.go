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
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/flow"
	"software_updater/core/hook"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util"
	"software_updater/core/util/error_util"
	"time"
)

type DefaultEngine struct {
	actionManager *ActionManager
	flowManager   *FlowManager
	taskRunner    *TaskRunner
	scheduler     Scheduler
	config        *config.EngineConfig
	driver        selenium.WebDriver
}

func (e *DefaultEngine) InitEngine(ctx context.Context, engineConfig *config.EngineConfig) error {
	e.actionManager = NewActionManager()
	e.flowManager = NewFlowManager()
	e.taskRunner = NewTaskRunner(ctx, true, engineConfig.RunnerCheck, func(ctx context.Context, cv *po.CurrentVersion, v *po.Version) {
		_ = e.updateCurrentVersion(ctx, v, cv)
	})
	e.scheduler = NewScheduler()
	e.config = engineConfig
	e.driver = web.Driver()
	return nil
}

func (e *DefaultEngine) DestroyEngine(ctx context.Context) {
	e.taskRunner.Stop(ctx)
}

func (e *DefaultEngine) RegisterAction(factory action.Factory) error {
	e.actionManager.Register(factory)
	return nil
}

func (e *DefaultEngine) RegisterHook(registerItem *hook.RegisterInfo) error {
	return e.actionManager.RegisterHook(registerItem)
}

func (e *DefaultEngine) Load(ctx context.Context, homepage *po.Homepage, useCache bool) (*flow.Flow, error) {
	return e.flowManager.Load(ctx, homepage.Name, homepage.Actions, e.actionManager, useCache)
}

func (e *DefaultEngine) Run(ctx context.Context, hp *po.Homepage) (TaskID, error) {
	fl, err := e.Load(ctx, hp, true)
	if err != nil {
		return 0, err
	}

	id, err := e.taskRunner.EnqueueJob(ctx, fl, hp.Name, hp.Current, hp.HomepageURL)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *DefaultEngine) CheckState(_ context.Context, id TaskID) (bool, flow.State, error) {
	return e.taskRunner.GetTaskState(id)
}

func (e *DefaultEngine) GetTaskIDMap(_ context.Context) (map[string]TaskID, error) {
	return e.taskRunner.GetTaskIDMap()
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
