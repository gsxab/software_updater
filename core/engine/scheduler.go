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
	"software_updater/core/db/po"
	"time"
)

type Scheduler interface {
	ScheduleForUpdate(cv *po.CurrentVersion, oldV *po.Version, newV *po.Version) time.Time
	ScheduleForNoUpdate(cv *po.CurrentVersion, oldV *po.Version) time.Time
}

func NewScheduler() Scheduler {
	return &DefaultScheduler{}
}

const day time.Duration = time.Hour * 24
const defaultSchedule time.Duration = day * 7

type DefaultScheduler struct {
}

func (s *DefaultScheduler) ScheduleForNoUpdate(_ *po.CurrentVersion, oldV *po.Version) time.Time {
	if oldV == nil {
		// first version not found, not expected
		return time.Now().Add(defaultSchedule)
	}

	var lastUpdate time.Time
	if oldV.RemoteDate != nil {
		lastUpdate = *oldV.RemoteDate
	} else if oldV.LocalTime != nil {
		lastUpdate = *oldV.LocalTime
	} else {
		// no available time for the previous version
		return time.Now().Add(defaultSchedule * 2)
	}

	thisUpdate := time.Now()

	daysToUpdate := thisUpdate.Sub(lastUpdate).Hours() / 24
	var daysToSchedule time.Duration
	if daysToUpdate < 10 {
		daysToSchedule = 7
	} else if daysToUpdate > 30 {
		daysToSchedule = 23
	} else {
		daysToSchedule = time.Duration(int(3 * daysToUpdate / 4))
	}
	nextUpdateTime := thisUpdate.Add(day * daysToSchedule)

	nextUpdateDate := time.Date(nextUpdateTime.Year(), nextUpdateTime.Month(), nextUpdateTime.Day(), 0, 0, 0, 0, time.Local)
	return nextUpdateDate
}

func (s *DefaultScheduler) ScheduleForUpdate(_ *po.CurrentVersion, oldV *po.Version, newV *po.Version) time.Time {
	if oldV == nil {
		// first fetched version
		return time.Now().Add(defaultSchedule)
	}

	var lastUpdate time.Time
	if oldV.RemoteDate != nil {
		lastUpdate = *oldV.RemoteDate
	} else if oldV.LocalTime != nil {
		lastUpdate = *oldV.LocalTime
	} else {
		// no available time for the previous version
		return time.Now().Add(defaultSchedule * 2)
	}

	var thisUpdate time.Time
	if newV.RemoteDate != nil {
		thisUpdate = *newV.RemoteDate
	} else if oldV.LocalTime != nil {
		thisUpdate = *newV.LocalTime
	} else {
		thisUpdate = time.Now()
	}

	daysToUpdate := thisUpdate.Sub(lastUpdate).Hours() / 24
	var daysToSchedule time.Duration
	if daysToUpdate < 14 {
		daysToSchedule = 7
	} else if daysToUpdate > 45 {
		daysToSchedule = 23
	} else {
		daysToSchedule = time.Duration(int(daysToUpdate / 2))
	}
	nextUpdateTime := thisUpdate.Add(day * daysToSchedule)

	nextUpdateDate := time.Date(nextUpdateTime.Year(), nextUpdateTime.Month(), nextUpdateTime.Day(), 0, 0, 0, 0, time.Local)
	return nextUpdateDate
}
