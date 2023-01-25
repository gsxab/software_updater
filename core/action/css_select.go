package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type CSSSelect struct {
	Default
	DefaultFactory[CSSSelect, *CSSSelect]
	Selector string `json:"selector"`
}

func (a *CSSSelect) Path() Path {
	return Path{"selector", "css", "css_select"}
}

func (a *CSSSelect) OutElmNum() int {
	return 1
}

func (a *CSSSelect) Do(ctx context.Context, driver selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = ElementToArgs(element, input)
	return
}

func (a *CSSSelect) ToDTO() *DTO {
	return &DTO{
		ReadPage: true,
		Output:   []string{"node"},
		Values:   map[string]string{"selector": a.Selector},
	}
}
