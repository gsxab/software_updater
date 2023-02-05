package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
	"time"
)

type WaitFor struct {
	base.Default
	base.DefaultFactory[WaitFor, *WaitFor]
	Delay string `json:"delay"`
	delay time.Duration
}

func (a *WaitFor) Path() action.Path {
	return action.Path{"basic", "wait", "delay"}
}

func (a *WaitFor) Icon() string {
	return "mdi:mdi-timer"
}

func (a *WaitFor) Init(ctx context.Context, _ *sync.WaitGroup) (err error) {
	a.delay, err = time.ParseDuration(a.Delay)
	if err != nil {
		logs.Error(ctx, "duration parsing failed", err, "duration", a.Delay)
		return
	}
	return
}

func (a *WaitFor) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = input
	t := time.NewTimer(a.delay)
	select {
	case <-ctx.Done():
		logs.WarnM(ctx, "delay action cancelled")
		exit = action.Cancelled
	case <-t.C:
	}
	if !t.Stop() {
		<-t.C
	}
	return
}

func (a *WaitFor) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"delay": a.Delay},
	}
}
