package po

import (
	"gorm.io/gorm"
	"time"
)

type CurrentVersion struct {
	gorm.Model
	Name             string     `gorm:"column:name;index;unique"`
	ScheduledAt      *time.Time `gorm:"column:scheduled_at"`
	CurrentVersionID uint       `gorm:"column:current_version_id;notNull"`
	Info             string     `gorm:"column:info"`
	CurrentVersion   *Version   `gorm:"references:CurrentVersionID;foreignKey:ID"`
}
