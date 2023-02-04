package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util/url_util"
	"sync"
)

type Access struct {
	base.Default
	base.DefaultFactory[Access, *Access]
}

func (a *Access) Path() action.Path {
	return action.Path{"browser", "access", "goto_url"}
}

func (a *Access) Icon() string {
	return "web"
}

func (a *Access) OutElmNum() int {
	return 0
}

func (a *Access) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	relURL := input.Strings[0]
	baseURL, err := driver.CurrentURL()
	if err != nil {
		logs.Error(ctx, "selenium current url failed", err)
		return
	}
	url, err := url_util.RelativeURL(baseURL, relURL)
	if err != nil {
		logs.Error(ctx, "relative url calculation failed", err, "baseURL", baseURL, "relURL", relURL)
		return
	}
	err = driver.Get(url.String())
	output = action.ElementsToArgs([]selenium.WebElement{}, input)
	return
}

func (a *Access) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			OpenPage: true,
			Input:    []string{"rel_url"},
		},
	}
}
