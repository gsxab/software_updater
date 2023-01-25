package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
)

type StoreDigest struct {
	Default
	DefaultFactory[StoreDigest, *StoreDigest]
	Index int `json:"index"`
}

func (a *StoreDigest) Path() Path {
	return Path{"basic", "value_store", "store_digest"}
}

func (a *StoreDigest) Do(ctx context.Context, _ selenium.WebDriver, input *Args, version *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "date storing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	version.Digest = &text
	output = input
	return
}

func (a *StoreDigest) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"text"},
		Values: map[string]string{"index": util.ToJSON(a.Index)},
	}
}
