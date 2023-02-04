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

type ReturnConst struct {
	base.Default
	base.DefaultFactory[ReturnConst, *ReturnConst]
	Val []string `json:"val"`
}

func (r *ReturnConst) Path() action.Path {
	return action.Path{"basic", "value_generator", "return_value"}
}

func (a *ReturnConst) Icon() string {
	return "text-box-outline"
}

func (a *ReturnConst) OutStrNum() int {
	return len(a.Val)
}

func (a *ReturnConst) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = action.StringsToArgs(a.Val, input)
	return
}

func (a *ReturnConst) ToDTO() *action.DTO {
	return &action.DTO{Values: map[string]string{"value": util.ToJSON(a.Val)}}
}
