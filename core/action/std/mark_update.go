package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"
)

type MarkUpdate struct {
	base.Default
}

func (a *MarkUpdate) Path() action.Path {
	return []string{"basic", "value_control", "mark_update"}
}

func (a *MarkUpdate) Icon() string {
	return "mdi:mdi-checkbox-marked-circle"
}

func (a *MarkUpdate) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return action.MarkUpdateToArgs(input), action.Finished, nil
}

func (a *MarkUpdate) ToDTO() *action.DTO {
	return &action.DTO{}
}

func (a *MarkUpdate) NewAction(_ string) (action.Action, error) {
	return a, nil
}

func (a *MarkUpdate) ToProtoDTO() *action.ProtoDTO {
	return nil
}
