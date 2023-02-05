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

type StoreURL struct {
	base.Default
	base.DefaultFactory[StoreURL, *StoreURL]
	base.IndexReader
}

func (a *StoreURL) Path() action.Path {
	return action.Path{"basic", "value_store", "store_url"}
}

func (a *StoreURL) Icon() string {
	return "mdi:mdi-text-box-check"
}

func (a *StoreURL) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Read(ctx, input, func(text string) {
		version.Link = &text
	})
}

func (a *StoreURL) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"text"},
		},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
