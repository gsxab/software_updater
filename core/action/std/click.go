package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type Click struct {
	prototype.Default
}

func (a *Click) Path() action.Path {
	return action.Path{"browser", "interact", "click"}
}

func (a *Click) InElmNum() int {
	return 1
}

func (a *Click) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = input
	err = input.Elements[0].Click()
	if err != nil {
		logs.Error(ctx, "selenium click failed", err)
		return
	}
	return
}

func (a *Click) ToDTO() *action.DTO {
	return &action.DTO{
		OpenPage: true,
		Input:    []string{"node"},
	}
}

func (a *Click) NewAction(string) (action.Action, error) {
	return a, nil
}
