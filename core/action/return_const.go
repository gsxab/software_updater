package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type ReturnConst struct {
	Default
	DefaultFactory[ReturnConst, *ReturnConst]
	Val []string `json:"val"`
}

func (r *ReturnConst) Path() Path {
	return Path{"basic", "value_generator", "return_value"}
}

func (a *ReturnConst) OutStrNum() int {
	return len(a.Val)
}

func (a *ReturnConst) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = StringsToArgs(a.Val, input)
	return
}

func (a *ReturnConst) ToDTO() *DTO {
	return &DTO{Values: map[string]string{"value": util.ToJSON(a.Val)}}
}
