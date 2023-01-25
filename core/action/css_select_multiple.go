package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type CSSSelectMultiple struct {
	Default
	DefaultFactory[CSSSelectMultiple, *CSSSelectMultiple]
	Selector string `json:"selector"`
}

func (a *CSSSelectMultiple) Path() Path {
	return Path{"selector", "css", "css_select_multiple"}
}

func (a *CSSSelectMultiple) OutElmNum() int {
	return Any
}

func (a *CSSSelectMultiple) Do(ctx context.Context, driver selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	elements, err := driver.FindElements(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = ElementsToArgs(elements, input)
	return
}

func (a *CSSSelectMultiple) ToDTO() *DTO {
	return &DTO{
		ReadPage: true,
		Output:   []string{"nodes..."},
		Values:   map[string]string{"selector": a.Selector},
	}
}
