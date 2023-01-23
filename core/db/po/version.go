package po

import (
	"gorm.io/gorm"
	"time"
)

type Version struct {
	gorm.Model
	Name       string          `gorm:"column:name;index:idx_versions_name_version"`
	Version    string          `gorm:"column:version;index:idx_versions_name_version"`
	Filename   *string         `gorm:"column:filename"`
	Picture    *string         `gorm:"column:picture"`
	Link       *string         `gorm:"column:link"`
	Digest     *string         `gorm:"column:digest"`
	RemoteDate *time.Time      `gorm:"column:remote_date;type:date"`
	LocalTime  *time.Time      `gorm:"column:local_time"`
	Previous   *string         `gorm:"column:previous_version"`
	CV         *CurrentVersion `gorm:"foreignKey:VersionID;references:ID"`
}
