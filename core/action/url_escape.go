package action

import (
	"context"
	"github.com/tebeka/selenium"
	"net/url"
	"software_updater/core/db/po"
	"sync"
)

type URLEscape struct {
	StringMutator
	DefaultFactory[URLEscape, *URLEscape]
}

func (a *URLEscape) Path() Path {
	return Path{"string", "mutator", "url_escape"}
}

func (a *URLEscape) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.Mutate(input, func(text string) string {
		return url.QueryEscape(text)
	})
}
