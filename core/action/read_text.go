package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type ReadText struct {
	Default
}

func (a *ReadText) Path() Path {
	return Path{"browser", "reader", "read_text"}
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

func (a *ReadText) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	text, err := input.Elements[0].Text()
	if err != nil {
		logs.Error(ctx, "selenium element get_text failed", err)
		return
	}
	output = StringToArgs(text, input)
	return
}

func (a *ReadText) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"node"},
		Output: []string{"text"},
	}
}

func (a *ReadText) NewAction(string) (Action, error) {
	return a, nil
}
