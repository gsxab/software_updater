package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util/url_util"
	"sync"
)

type Access struct {
	Default
	DefaultFactory[Access, *Access]
}

func (a *Access) Path() Path {
	return Path{"browser", "access", "goto_url"}
}

func (a *Access) OutElmNum() int {
	return 0
}

func (a *Access) Do(ctx context.Context, driver selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
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
	output = ElementsToArgs([]selenium.WebElement{}, input)
	return
}

func (a *Access) ToDTO() *DTO {
	return &DTO{
		OpenPage: true,
		Input:    []string{"rel_url"},
	}
}
