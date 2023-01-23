package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type StoreURL struct {
	Default
	DefaultFactory[StoreURL, *StoreURL]
	Index int `json:"index"`
}

func (a *StoreURL) Path() Path {
	return Path{"basic", "value_store", "store_url"}
}

func (a *StoreURL) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		return
	}
	text := input.Strings[a.Index]
	version.Link = &text
	output = input
	return
}

func (a *StoreURL) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"text"},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
