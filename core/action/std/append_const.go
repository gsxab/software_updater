package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type AppendConst struct {
	base.Default
	base.DefaultFactory[AppendConst, *AppendConst]
	Val string `json:"val"`
}

func (a *AppendConst) Path() action.Path {
	return action.Path{"basic", "value_generator", "append_value"}
}

func (a *AppendConst) Icon() string {
	return "mdi:mdi-text-box-plus"
}

func (a *AppendConst) OutStrNum() int {
	return len(a.Val)
}

func (a *AppendConst) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = action.AnotherStringToArgs(a.Val, input)
	return
}

func (a *AppendConst) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"value": util.ToJSON(a.Val)},
	}
}
