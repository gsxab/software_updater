package std

import (
	"context"
	"github.com/tebeka/selenium"
	"net/url"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"sync"
)

type URLEscape struct {
	prototype.StringMutator
	prototype.DefaultFactory[URLEscape, *URLEscape]
}

func (a *URLEscape) Path() action.Path {
	return action.Path{"string", "mutator", "url_escape"}
}

func (a *URLEscape) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Mutate(input, func(text string) string {
		return url.QueryEscape(text)
	})
}
