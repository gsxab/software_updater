package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type Click struct {
	base.Default
}

func (a *Click) Path() action.Path {
	return action.Path{"browser", "interact", "click"}
}

func (a *Click) Icon() string {
	return "mdi:mdi-cursor-default-click"
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
		ProtoDTO: a.ToProtoDTO(),
	}
}

func (a *Click) NewAction(string) (action.Action, error) {
	return a, nil
}

func (a *Click) ToProtoDTO() *action.ProtoDTO {
	return &action.ProtoDTO{
		OpenPage: true,
		Input:    []string{"node"},
	}
}
