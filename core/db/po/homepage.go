package po

import "gorm.io/gorm"

type Homepage struct {
	gorm.Model
	Name        string          `gorm:"column:name;notNull"`
	HomepageURL string          `gorm:"column:homepage_url;notNull"`
	Actions     string          `gorm:"column:version_actions;notNull;default:{}"`
	NoUpdate    bool            `gorm:"column:no_update"`
	Current     *CurrentVersion `gorm:"references:Name;foreignKey:Name"`
}
