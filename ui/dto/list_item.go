package dto

import (
	"software_updater/core/db/po"
	"software_updater/core/util"
	"time"
)

type ListItemDTO struct {
	Name       string  `json:"name"`
	Version    *string `json:"version"`
	UpdateDate *string `json:"update_date"`
	SchedDate  *string `json:"scheduled_date"`
}

func NewListItemDTO(homepage *po.Homepage, dateFormat string) *ListItemDTO {
	var updateAt, scheduledAt *time.Time
	var version *string
	if homepage.Current != nil {
		scheduledAt = homepage.Current.ScheduledAt
		if homepage.Current.Version != nil {
			updateAt = homepage.Current.Version.LocalTime
			version = &homepage.Current.Version.Version
		}
	}

	return &ListItemDTO{
		Name:       homepage.Name,
		Version:    version,
		UpdateDate: util.FormatTime(updateAt, dateFormat),
		SchedDate:  util.FormatTime(scheduledAt, dateFormat),
	}
}
