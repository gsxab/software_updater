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

type StoreVersion struct {
	base.Default
	base.DefaultFactory[StoreVersion, *StoreVersion]
	base.IndexReader
}

func (a *StoreVersion) Path() action.Path {
	return action.Path{"basic", "value_store", "store_version"}
}

func (a *StoreVersion) Icon() string {
	return "text-box-edit-outline"
}

func (a *StoreVersion) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Read(ctx, input, func(text string) {
		version.Version = text
	})
}

func (a *StoreVersion) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"text"},
		},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}