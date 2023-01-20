package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
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

func (a *CSSSelect) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
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
