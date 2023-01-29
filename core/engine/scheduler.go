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
	return thisUpdate.Add(day * daysToSchedule)
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
	return thisUpdate.Add(day * daysToSchedule)
}
