package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type ReturnEmpty struct {
	Default
}

func (a *ReturnEmpty) Path() Path {
	return Path{"basic", " value_gen", " return_empty"}
}

func (a *ReturnEmpty) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return
}

func (a *ReturnEmpty) ToDTO() *DTO {
	return &DTO{}
}

func (a *ReturnEmpty) NewAction(_ string) (Action, error) {
	return a, nil
}
