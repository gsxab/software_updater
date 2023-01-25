package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
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
	return Path{"basic", "wait", "delay"}
}

func (a *WaitFor) Init(ctx context.Context, _ *sync.WaitGroup) (err error) {
	a.delay, err = time.ParseDuration(a.Delay)
	if err != nil {
		logs.Error(ctx, "duration parsing failed", err, "duration", a.Delay)
		return
	}
	return
}

func (a *WaitFor) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = input
	t := time.NewTimer(a.delay)
	select {
	case <-ctx.Done():
		logs.WarnM(ctx, "delay action cancelled")
		exit = Cancelled
	case <-t.C:
	}
	if !t.Stop() {
		<-t.C
	}
	return
}

func (a *WaitFor) ToDTO() *DTO {
	return &DTO{
		Values: map[string]string{"delay": a.Delay},
	}
}
