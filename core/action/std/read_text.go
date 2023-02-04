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

type ReadText struct {
	prototype.Default
}

func (a *ReadText) Path() action.Path {
	return action.Path{"browser", "reader", "read_text"}
}

func (a *ReadText) InElmNum() int {
	return 1
}

func (a *ReadText) OutElmNum() int {
	return 1
}

func (a *ReadText) OutStrNum() int {
	return 1
}

func (a *ReadText) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	text, err := input.Elements[0].Text()
	if err != nil {
		logs.Error(ctx, "selenium element get_text failed", err)
		return
	}
	output = action.StringToArgs(text, input)
	return
}

func (a *ReadText) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"node"},
		Output: []string{"text"},
	}
}

func (a *ReadText) NewAction(string) (action.Action, error) {
	return a, nil
}
