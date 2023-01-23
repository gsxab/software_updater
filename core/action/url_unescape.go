package action

import (
	"context"
	"github.com/tebeka/selenium"
	"net/url"
	"software_updater/core/db/po"
	"sync"
)

type URLUnescape struct {
	StringMutator
	DefaultFactory[URLUnescape, *URLUnescape]
}

func (a *URLUnescape) Path() Path {
	return Path{"string", "mutator", "url_unescape"}
}

func (a *URLUnescape) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(input, func(text string) (string, error) {
		return url.QueryUnescape(text)
	})
}
