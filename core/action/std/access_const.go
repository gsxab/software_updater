package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util/url_util"
	"sync"
)

type AccessConst struct {
	base.Default
	base.DefaultFactory[AccessConst, *AccessConst]
	URL string `json:"url"`
}

func (a *AccessConst) Path() action.Path {
	return action.Path{"browser", "access", "goto_const"}
}

func (a *AccessConst) Icon() string {
	return "mdi:mdi-web"
}

func (a *AccessConst) OutElmNum() int {
	return 0
}

func (a *AccessConst) Do(ctx context.Context, driver selenium.WebDriver, _ *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	base, err := driver.CurrentURL()
	if err != nil {
		logs.Error(ctx, "selenium get current url failed", err)
		return
	}
	url, err := url_util.RelativeURL(base, a.URL)
	if err != nil {
		logs.Error(ctx, "relative url calculation failed", err)
		return
	}
	err = driver.Get(url.String())
	if err != nil {
		logs.Error(ctx, "selenium url access failed", err, "URL", a.URL)
		return
	}
	size := config.Current().Selenium.WindowSize
	err = driver.ResizeWindow("", size.Width, size.Height)
	if err != nil {
		logs.Error(ctx, "selenium resize failed", err)
		return
	}
	return
}

func (a *AccessConst) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			OpenPage: true,
		},
		Values: map[string]string{"url": a.URL},
	}
}
