package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
	"time"
)

type StoreDate struct {
	Default
	DefaultFactory[StoreDate, *StoreDate]
	Index  int    `json:"index"`
	Format string `json:"format"`
}

func (a *StoreDate) Path() Path {
	return Path{"basic", "value_store", "store_date"}
}

func (a *StoreDate) Do(ctx context.Context, _ selenium.WebDriver, input *Args, version *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "date storing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	t, err := time.Parse(a.Format, text) // ignore tz. only date to be shown, and we don't know the right time zones.
	if err != nil {
		logs.Error(ctx, "date parsing failed", err, "text", text)
		return
	}
	version.RemoteDate = &t
	output = input
	return
}

func (a *StoreDate) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"text"},
		Values: map[string]string{"index": util.ToJSON(a.Index), "format": a.Format},
	}
}
