package ui

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/job"
	"software_updater/core/util"
	"software_updater/ui/dto"
)

type App struct {
	current    *po.Homepage
	window     fyne.Window
	rootTab    *container.AppTabs
	listData   []*dto.ListItemDTO
	listPage   fyne.CanvasObject
	detailData []*dto.VersionDTO
	detailPage fyne.CanvasObject
	flowData   *job.Flow
	flowPage   fyne.CanvasObject
}

func (a *App) initGUI(ctx context.Context) error {
	err := a.reloadData(ctx)
	if err != nil {
		return err
	}
	listTable := widget.NewTable(
		func() (int, int) {
			return len(a.listData) + 1, 4
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("stub"))
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			if id.Row == 0 {
				headers := []string{"name", "version", "update date", "scheduled date", "actions"}
				label := object.(*fyne.Container).Objects[0].(*widget.Label)
				label.SetText(headers[id.Col])
				return
			}
			switch id.Col {
			case 0:
				label := object.(*fyne.Container).Objects[0].(*widget.Label)
				label.SetText(a.listData[id.Row-1].Name)
			case 1:
				label := object.(*fyne.Container).Objects[0].(*widget.Label)
				label.SetText(util.Default(a.listData[id.Row-1].Version, ""))
			case 2:
				label := object.(*fyne.Container).Objects[0].(*widget.Label)
				label.SetText(util.Default(a.listData[id.Row-1].UpdateDate, ""))
			case 3:
				label := object.(*fyne.Container).Objects[0].(*widget.Label)
				label.SetText(util.Default(a.listData[id.Row-1].SchedDate, ""))
			}
		},
	)
	columnWidthList := []float32{150, 100, 150, 150}
	for idx, listWidth := range columnWidthList {
		listTable.SetColumnWidth(idx, listWidth)
	}

	detailTable := widget.NewTable(
		func() (int, int) {
			return len(a.detailData), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("stub")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
		},
	)
	detailLast := widget.NewToolbarAction(theme.NavigateBackIcon(), func() {})
	detailLastSep := widget.NewToolbarSeparator()
	detailDownload := widget.NewToolbarAction(theme.DownloadIcon(), func() {})
	detailNextSep := widget.NewToolbarSeparator()
	detailNext := widget.NewToolbarAction(theme.NavigateNextIcon(), func() {})
	detailToolbar := widget.NewToolbar(detailLast, detailLastSep, detailDownload, detailNextSep, detailNext)
	detailPage := container.NewVBox(detailToolbar, detailTable)

	a.listPage = listTable
	a.detailPage = detailPage
	a.flowPage = canvas.NewText("Flow page", nil)

	a.rootTab = container.NewAppTabs(
		container.NewTabItemWithIcon("List", theme.ListIcon(), a.listPage),
		container.NewTabItemWithIcon("Detail", theme.InfoIcon(), a.detailPage),
		container.NewTabItemWithIcon("Flow", theme.FileIcon(), a.flowPage),
	)
	a.rootTab.SetTabLocation(container.TabLocationLeading)
	a.rootTab.DisableIndex(1)
	a.rootTab.DisableIndex(2)
	a.window.SetContent(a.rootTab)
	return nil
}

func (a *App) reloadData(ctx context.Context) error {
	data := make([]*dto.ListItemDTO, 0)
	homepages, err := dao.Homepage.WithContext(ctx).LeftJoin(dao.CurrentVersion).LeftJoin(dao.Version).Find()
	if err != nil {
		return err
	}
	for _, homepage := range homepages {
		data = append(data, dto.NewListItemDTO(homepage, uiConfig.DateFormat))
	}
	a.listData = data
	return nil
}

func InitAndRun(ctx context.Context, configExtraUI string) error {
	var err error

	uiConfig = DefaultConfig()
	if configExtraUI != "" {
		err = json.Unmarshal([]byte(configExtraUI), uiConfig)
		if err != nil {
			return err
		}
	}

	a := app.New()
	myApp := App{}
	myApp.window = a.NewWindow("Software Updater")

	myApp.window.Resize(fyne.Size{
		Width:  640,
		Height: 480,
	})
	err = myApp.initGUI(ctx)
	if err != nil {
		return err
	}

	myApp.window.ShowAndRun()
	return nil
}
