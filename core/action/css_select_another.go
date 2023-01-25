package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type CSSSelectAppend struct {
	Default
	DefaultFactory[CSSSelectAppend, *CSSSelectAppend]
	Selector string `json:"selector"`
}

func (a *CSSSelectAppend) Path() Path {
	return Path{"selector", "css", "css_select_another"}
}

func (a *CSSSelectAppend) OutElmNum() int {
	return OneMore
}

func (a *CSSSelectAppend) Do(ctx context.Context, driver selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = AnotherElementToArgs(element, input)
	return
}

func (a *CSSSelectAppend) ToDTO() *DTO {
	return &DTO{
		ReadPage: true,
		Input:    []string{"nodes..."},
		Output:   []string{"node", "nodes..."},
		Values:   map[string]string{"selector": a.Selector},
	}
}
