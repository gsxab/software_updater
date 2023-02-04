package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type StoreDigest struct {
	prototype.Default
	prototype.DefaultFactory[StoreDigest, *StoreDigest]
	prototype.IndexReader
}

func (a *StoreDigest) Path() action.Path {
	return action.Path{"basic", "value_store", "store_digest"}
}

func (a *StoreDigest) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Read(ctx, input, func(text string) {
		version.Digest = &text
	})
}

func (a *StoreDigest) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"text"},
		Output: []string{"text"},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
