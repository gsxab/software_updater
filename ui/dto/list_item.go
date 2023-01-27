package dto

type ListItemDTO struct {
	Name          string  `json:"name" gorm:"column:name"`
	PageURL       string  `json:"page_url" gorm:"column:homepage_url"`
	Version       *string `json:"version" gorm:"column:version"`
	UpdateDate    *string `json:"update_date" gorm:"column:local_time"`
	ScheduledDate *string `json:"scheduled_date" gorm:"column:scheduled_at"`
}
