package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
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

func (a *CSSSelectAppend) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	element, err := driver.FindElement(selenium.ByCSSSelector, a.Selector)
	if err != nil {
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
