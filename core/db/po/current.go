package po

import (
	"gorm.io/gorm"
	"time"
)

type CurrentVersion struct {
	gorm.Model
	Name        string     `gorm:"column:name;index;unique"`
	ScheduledAt *time.Time `gorm:"column:scheduled_at"`
	VersionID   uint       `gorm:"column:version_id;notNull"`
	Info        string     `gorm:"column:info"`
	Version     *Version   `gorm:"joinForeignKey:ID;joinReferences:VersionID"`
}
