package po

import "gorm.io/gorm"

type Homepage struct {
	gorm.Model
	Name        string          `gorm:"column:name;notNull"`
	HomepageURL string          `gorm:"column:homepage_url;notNull"`
	Actions     string          `gorm:"column:version_actions;notNull"`
	NoUpdate    bool            `gorm:"column:update"`
	Current     *CurrentVersion `gorm:"references:Name;foreignKey:Name"`
}
