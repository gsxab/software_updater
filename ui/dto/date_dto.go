package dto

import (
	"math"
	"time"
)

type DateDTO struct {
	TS   int64 `json:"ts"`
	Days int64 `json:"days"`
}

func ToDateDTO(t *time.Time, loc *time.Location) *DateDTO {
	if t == nil {
		return nil
	}
	local := t.In(loc) // date saved in localtime zone
	date := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
	now := time.Now()
	days := int64(math.Floor(now.Sub(date).Hours() / 24))
	return &DateDTO{
		TS:   t.UnixMilli(),
		Days: days,
	}
}
