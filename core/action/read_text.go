package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
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

func (a *ReadText) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	text, err := input.Elements[0].Text()
	if err != nil {
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
