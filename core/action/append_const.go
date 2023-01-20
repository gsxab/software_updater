package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type AppendConst struct {
	Default
	DefaultFactory[AppendConst, *AppendConst]
	Val string `json:"val"`
}

func (a *AppendConst) Path() Path {
	return Path{"basic", "value_generator", "append_value"}
}

func (a *AppendConst) OutStrNum() int {
	return len(a.Val)
}

func (a *AppendConst) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = AnotherStringToArgs(a.Val, input)
	return
}

func (a *AppendConst) ToDTO() *DTO {
	return &DTO{Values: map[string]string{"value": util.ToJSON(a.Val)}}
}
