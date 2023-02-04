package std_action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type AccessConst struct {
	action.Default
	action.DefaultFactory[AccessConst, *AccessConst]
	URL string `json:"url"`
}

func (a *AccessConst) Path() action.Path {
	return action.Path{"browser", "access", "goto_url"}
}

func (a *AccessConst) OutElmNum() int {
	return 0
}

func (a *AccessConst) Do(ctx context.Context, driver selenium.WebDriver, _ *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	err = driver.Get(a.URL)
	if err != nil {
		logs.Error(ctx, "selenium url access failed", err, "URL", a.URL)
		return
	}
	return
}

func (a *AccessConst) ToDTO() *action.DTO {
	return &action.DTO{
		OpenPage: true,
		Values:   map[string]string{"url": a.URL},
	}
}
