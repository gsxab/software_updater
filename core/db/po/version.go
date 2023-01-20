package po

import (
	"gorm.io/gorm"
	"time"
)

type Version struct {
	gorm.Model
	Name       string     `gorm:"column:name;index:name_version"`
	Version    string     `gorm:"column:version;index:name_version"`
	Filename   *string    `gorm:"column:filename"`
	Picture    *string    `gorm:"column:picture"`
	Link       *string    `gorm:"column:link"`
	Digest     *string    `gorm:"column:digest"`
	RemoteDate *time.Time `gorm:"column:remote_update;type:date"`
	LocalTime  *time.Time `gorm:"column:local_time;type:date"`
}
