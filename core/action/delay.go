package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
	"time"
)

type WaitFor struct {
	Default
	DefaultFactory[WaitFor, *WaitFor]
	Delay string `json:"delay"`
	delay time.Duration
}

func (a *WaitFor) Path() Path {
	return Path{"basic", " wait", " delay"}
}

func (a *WaitFor) Init(context.Context, *sync.WaitGroup) (err error) {
	a.delay, err = time.ParseDuration(a.Delay)
	return
}

func (a *WaitFor) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = input
	select {
	case <-ctx.Done():
		exit = Cancelled
	case <-time.After(a.delay):
	}
	return
}

func (a *WaitFor) ToDTO() *DTO {
	return &DTO{
		Values: map[string]string{"delay": a.Delay},
	}
}
