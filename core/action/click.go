package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type Click struct {
	Default
}

func (a *Click) Path() Path {
	return Path{"browser", "interact", "click"}
}

func (a *Click) InElmNum() int {
	return 1
}

func (a *Click) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = input
	err = input.Elements[0].Click()
	if err != nil {
		return
	}
	return
}

func (a *Click) ToDTO() *DTO {
	return &DTO{
		OpenPage: true,
		Input:    []string{"node"},
	}
}

func (a *Click) NewAction(string) (Action, error) {
	return a, nil
}
