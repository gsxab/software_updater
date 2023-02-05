package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
	"time"
)

type StoreDate struct {
	base.Default
	base.DefaultFactory[StoreDate, *StoreDate]
	base.IndexReader
	Format string `json:"format"`
}

func (a *StoreDate) Path() action.Path {
	return action.Path{"basic", "value_store", "store_date"}
}

func (a *StoreDate) Icon() string {
	return "mdi:mdi-clock-check"
}

func (a *StoreDate) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.ReadWithErr(ctx, input, func(text string) error {
		t, err := time.Parse(a.Format, text) // implicitly use UTC. only date to be shown, and we don't know the right time zones.
		if err != nil {
			logs.Error(ctx, "date parsing failed", err, "text", text)
			return err
		}
		version.RemoteDate = &t
		return nil
	})
}

func (a *StoreDate) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"text"},
		},
		Values: map[string]string{"index": util.ToJSON(a.Index), "format": a.Format},
	}
}
