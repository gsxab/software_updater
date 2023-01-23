package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type StoreVersion struct {
	Default
	DefaultFactory[StoreVersion, *StoreVersion]
	Index int `json:"index"`
}

func (a *StoreVersion) Path() Path {
	return Path{"basic", "value_store", "store_version"}
}

func (a *StoreVersion) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		return
	}
	text := input.Strings[a.Index]
	version.Version = text
	output = input
	return
}

func (a *StoreVersion) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"text"},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
