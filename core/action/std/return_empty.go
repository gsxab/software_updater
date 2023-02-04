package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"sync"
)

type ReturnEmpty struct {
	prototype.Default
}

func (a *ReturnEmpty) Path() action.Path {
	return action.Path{"basic", "value_gen", "return_empty"}
}

func (a *ReturnEmpty) Do(context.Context, selenium.WebDriver, *action.Args, *po.Version, *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return
}

func (a *ReturnEmpty) ToDTO() *action.DTO {
	return &action.DTO{}
}

func (a *ReturnEmpty) NewAction(_ string) (action.Action, error) {
	return a, nil
}
