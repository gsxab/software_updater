package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"software_updater/ui/dto"
)

type ListItemUI struct {
	*fyne.Container
	Name     *widget.Label
	Version  *widget.Label
	Update   *widget.Label
	Schedule *widget.Label
	Data     *dto.ListItemDTO
}

func (l *ListItemUI) Rebind(dto *dto.ListItemDTO) {
	l.Data = dto
	l.Name.SetText(dto.Name)
}

func NewListItemUI() *ListItemUI {
	l := &ListItemUI{
		Data: &dto.ListItemDTO{
			Name:          "",
			Version:       nil,
			UpdateDate:    nil,
			ScheduledDate: nil,
		},
	}

	l.Container = container.NewHBox(l.Name, l.Version, l.Update, l.Schedule)
	return l
}
