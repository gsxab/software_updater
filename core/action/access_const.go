package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type AccessConst struct {
	Default
	DefaultFactory[AccessConst, *AccessConst]
	URL string `json:"url"`
}

func (a *AccessConst) Path() Path {
	return Path{"browser", "access", "goto_url"}
}

func (a *AccessConst) OutElmNum() int {
	return 0
}

func (a *AccessConst) Do(ctx context.Context, driver selenium.WebDriver, _ *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	err = driver.Get(a.URL)
	if err != nil {
		logs.Error(ctx, "selenium url access failed", err, "URL", a.URL)
		return
	}
	return
}

func (a *AccessConst) ToDTO() *DTO {
	return &DTO{
		OpenPage: true,
		Values:   map[string]string{"url": a.URL},
	}
}
