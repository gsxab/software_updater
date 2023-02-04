package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type CSSSelectMultiple struct {
	prototype.Default
	prototype.DefaultFactory[CSSSelectMultiple, *CSSSelectMultiple]
	Selector string `json:"selector"`
}

func (a *CSSSelectMultiple) Path() action.Path {
	return action.Path{"selector", "css", "css_select_multiple"}
}

func (a *CSSSelectMultiple) OutElmNum() int {
	return action.Any
}

func (a *CSSSelectMultiple) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements, err := driver.FindElements(selenium.ByCSSSelector, a.Selector)
	if err != nil {
		logs.Error(ctx, "selenium element finding failed", err, "selector", a.Selector)
		return
	}
	output = action.ElementsToArgs(elements, input)
	return
}

func (a *CSSSelectMultiple) ToDTO() *action.DTO {
	return &action.DTO{
		ReadPage: true,
		Output:   []string{"nodes..."},
		Values:   map[string]string{"selector": a.Selector},
	}
}
