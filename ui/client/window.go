package client

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image"
	"net/url"
	"os"
	"path"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/job"
	"software_updater/core/logs"
	"software_updater/core/util/optional"
	"software_updater/ui/common"
	"software_updater/ui/dto"
)

type App struct {
	current *po.Homepage
	window  fyne.Window
	rootTab *container.AppTabs

	// list tab
	listData []*dto.ListItemDTO
	listPage fyne.CanvasObject

	// detail tab
	detailData     *dto.VersionDTO
	detailVal      [][]string
	detailPic      *canvas.Image
	detailPage     fyne.CanvasObject
	detailName     *widget.Label
	detailVer      *widget.Label
	detailList     *widget.Label
	detailToolLast *widget.ToolbarAction
	detailToolNext *widget.ToolbarAction

	flowData *job.Flow
	flowPage fyne.CanvasObject
}

func (a *App) initGUI(ctx context.Context, fa fyne.App) error {
	err := a.reloadListData(ctx)
	if err != nil {
		return err
	}

	// list tab
	listTable := widget.NewList(
		func() int {
			return len(a.listData)
		},
		func() fyne.CanvasObject {
			name := widget.NewLabel("loading...")
			version := widget.NewLabel("")
			update := widget.NewLabel("")
			schedule := widget.NewLabel("")

			detail := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {})
			run := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
			del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})
			del.Disable()
			actions := container.NewHBox(detail, run, del)

			hBox := container.NewGridWithColumns(5, name, version, update, schedule, actions)
			return hBox
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			objects := object.(*fyne.Container).Objects
			objects[0].(*widget.Label).SetText(a.listData[id].Name)
			objects[1].(*widget.Label).SetText(optional.ValueOr(a.listData[id].Version, "(no information available)"))
			objects[2].(*widget.Label).SetText(optional.ValueOr(a.listData[id].UpdateDate, "(not updated)"))
			objects[3].(*widget.Label).SetText(optional.ValueOr(a.listData[id].ScheduledDate, "(not scheduled)"))

			buttons := objects[4].(*fyne.Container).Objects
			if a.listData[id].Version != nil {
				buttons[0].(*widget.Button).Enable()
			} else {
				buttons[0].(*widget.Button).Disable()
			}
			buttons[0].(*widget.Button).OnTapped = func() {
				err = a.selectDetailVersion(ctx, id)
				if err != nil {
					return
				}
				a.rootTab.SelectIndex(1)
			}
			buttons[1].(*widget.Button).OnTapped = func() { /*a.RunCrawl(id)*/ }
			buttons[2].(*widget.Button).OnTapped = func() { /*a.Delete(id)*/ }
		},
	)
	// list tool
	listPage := listTable

	// detail tab
	a.detailName = widget.NewLabel("loading...")
	a.detailVer = widget.NewLabel("")
	detailTitle := container.NewCenter(container.NewHBox(a.detailName, a.detailVer))
	detailList := widget.NewList(
		func() int {
			return len(a.detailVal)
		},
		func() fyne.CanvasObject {
			left := widget.NewLabel("")
			right := container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.SearchIcon(), func() {}), widget.NewLabel(""))
			return container.NewGridWithColumns(2, left, right)
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			fieldName := object.(*fyne.Container).Objects[0].(*widget.Label)
			fieldName.SetText(a.detailVal[id][0])
			fieldName.Wrapping = fyne.TextWrapBreak
			fieldVal := object.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
			fieldVal.SetText(a.detailVal[id][1])
			fieldVal.Wrapping = fyne.TextWrapBreak
			fieldButton := object.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
			if len(a.detailVal[id]) > 2 {
				if valType, valTarget := a.detailVal[id][2], a.detailVal[id][3]; len(valType) > 0 && len(valTarget) > 0 {
					switch valType {
					case "URL":
						fieldButton.Icon = theme.SearchIcon()
						fieldButton.OnTapped = func() {
							targetURL, err := url.Parse(valTarget)
							if err != nil {
								logs.Error(ctx, "parse url failed", err, "url", valTarget)
								return
							}
							err = fa.OpenURL(targetURL)
							if err != nil {
								logs.Error(ctx, "open url failed", err, "url", valTarget)
								return
							}
						}
					}
					fieldButton.Show()
				} else {
					fieldButton.Hide()
				}
			} else {
				fieldButton.Hide()
			}
		},
	)
	a.detailPic = canvas.NewImageFromResource(theme.FileImageIcon())
	a.detailPic.FillMode = canvas.ImageFillContain
	// detail tool
	a.detailToolLast = widget.NewToolbarAction(theme.NavigateBackIcon(), func() {})
	a.detailToolNext = widget.NewToolbarAction(theme.NavigateNextIcon(), func() {})
	detailToolRefresh := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {})
	detailToolDownload := widget.NewToolbarAction(theme.DownloadIcon(), func() {})
	detailToolEdit := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {})
	detailToolbar := widget.NewToolbar(
		a.detailToolLast, a.detailToolNext,
		widget.NewToolbarSeparator(),
		detailToolRefresh, detailToolDownload, detailToolEdit,
	)
	columns := container.NewHSplit(detailList, container.NewScroll(a.detailPic))
	columns.SetOffset(0.8)
	detailPage := container.NewBorder(container.NewVBox(detailToolbar, detailTitle), nil, nil, nil, columns)

	a.listPage = listPage
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

func (a *App) reloadListData(ctx context.Context) (err error) {
	a.listData, err = common.GetList(ctx, clientUIConfig.DateFormat)
	return err
}

func (a *App) selectDetailVersion(ctx context.Context, id int) error {
	return a.reloadDetailVersion(ctx, id, *a.listData[id].Version)
}

func (a *App) reloadDetailVersion(ctx context.Context, id int, v string) error {
	a.rootTab.DisableIndex(1)
	a.rootTab.DisableIndex(2)

	name := a.listData[id].Name
	page := a.listData[id].PageURL

	version, err := common.GetVersionDetail(ctx, name, &page, v, clientUIConfig.DateFormat)
	if err != nil {
		return err
	}

	a.rootTab.EnableIndex(1)
	if version.Picture != nil {
		filename := *version.Picture
		pathname := path.Join(config.Current().Files.ScreenshotDir, filename)
		file, err := os.Open(pathname)
		if err != nil {
			logs.Error(ctx, "file opening failed", err, "filename", filename, "pathname", pathname)
			return err
		}
		img, _, err := image.Decode(file)
		if err != nil {
			logs.Error(ctx, "image decoding failed", err, "filename", filename, "pathname", pathname)
			return err
		}
		a.detailPic.Resource, err = fyne.LoadResourceFromPath(pathname)
		if err != nil {
			logs.Error(ctx, "image loading failed", err, "filename", filename, "pathname", pathname)
			return err
		}
		width := float32(480) // a.detailPic.MinSize().Width
		a.detailPic.SetMinSize(fyne.Size{Width: width, Height: float32(img.Bounds().Dy()) * width / float32(img.Bounds().Dx())})
	}
	if version.PrevVersion != nil {
		a.detailToolNext.OnActivated = func() {
			err = a.reloadDetailVersion(ctx, id, *version.PrevVersion)
			if err != nil {
				logs.ErrorE(ctx, err, "id", id, "version", *version.PrevVersion)
			}
		}
	} else {
		a.detailToolNext.OnActivated = func() {}
	}
	if version.NextVersion != nil {
		a.detailToolNext.OnActivated = func() {
			err = a.reloadDetailVersion(ctx, id, *version.NextVersion)
			if err != nil {
				logs.ErrorE(ctx, err, "id", id, "version", *version.NextVersion)
			}
		}
	} else {
		a.detailToolNext.OnActivated = func() {}
	}

	a.detailName.SetText(name)
	a.detailVer.SetText(version.Version)
	a.fillDetailVal()

	return nil
}

func (a *App) fillDetailVal() {
	a.detailVal = make([][]string, 0, 32)
	a.detailVal = append(a.detailVal, []string{"page", a.detailData.PageURL, "URL", a.detailData.PageURL})
	a.detailVal = append(a.detailVal, []string{"version", a.detailData.Version})
	if a.detailData.PrevVersion != nil {
		a.detailVal = append(a.detailVal, []string{"previous version", *a.detailData.PrevVersion})
	}
	if a.detailData.NextVersion != nil {
		a.detailVal = append(a.detailVal, []string{"next version", *a.detailData.NextVersion})
	}
	if a.detailData.RemoteDate != nil {
		a.detailVal = append(a.detailVal, []string{"remote date", *a.detailData.RemoteDate})
	}
	a.detailVal = append(a.detailVal, []string{"update date", a.detailData.UpdateDate})
	if a.detailData.Link != nil {
		a.detailVal = append(a.detailVal, []string{"link", *a.detailData.Link, "URL", *a.detailData.Link})
	}
	if a.detailData.Digest != nil {
		a.detailVal = append(a.detailVal, []string{"digest", *a.detailData.Digest})
	}
}

func InitAndRun(ctx context.Context, configExtraUI string) error {
	var err error

	clientUIConfig = DefaultConfig()
	if configExtraUI != "" {
		err = json.Unmarshal([]byte(configExtraUI), clientUIConfig)
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
	err = myApp.initGUI(ctx, a)
	if err != nil {
		return err
	}

	myApp.window.ShowAndRun()
	return nil
}
