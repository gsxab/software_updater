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

func (a *URLUnescape) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		return url.QueryUnescape(text)
	})
}
