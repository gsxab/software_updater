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

type CSSSelectAppend struct {
	base.Default
	base.DefaultFactory[CSSSelectAppend, *CSSSelectAppend]
	Selector string `json:"selector"`
}

func (a *CSSSelectAppend) Path() action.Path {
	return action.Path{"browser", "selector", "css", "css_select_another"}
}

func (a *CSSSelectAppend) Icon() string {
	return "mdi:mdi-select-drag"
}

func (a *CSSSelectAppend) OutElmNum() int {
	return action.OneMore
}

func (a *CSSSelectAppend) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = action.AnotherElementToArgs(element, input)
	return
}

func (a *CSSSelectAppend) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			ReadPage: true,
			Input:    []string{"nodes..."},
			Output:   []string{"node", "nodes..."},
		},
		Values: map[string]string{"selector": a.Selector},
	}
}
