package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type CSSSelect struct {
	base.Default
	base.DefaultFactory[CSSSelect, *CSSSelect]
	Selector string `json:"selector"`
}

func (a *CSSSelect) Path() action.Path {
	return action.Path{"browser", "selector", "css", "css_select"}
}

func (a *CSSSelect) Icon() string {
	return "mdi:mdi-select"
}

func (a *CSSSelect) OutElmNum() int {
	return 1
}

func (a *CSSSelect) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = action.ElementToArgs(element, input)
	return
}

func (a *CSSSelect) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			ReadPage: true,
			Output:   []string{"node"},
		},
		Values: map[string]string{"selector": a.Selector},
	}
}
