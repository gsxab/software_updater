package std_action

import (
	"context"
	"github.com/tebeka/selenium"
	"net/url"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"sync"
)

type URLUnescape struct {
	action.StringMutator
	action.DefaultFactory[URLUnescape, *URLUnescape]
}

func (a *URLUnescape) Path() action.Path {
	return action.Path{"string", "mutator", "url_unescape"}
}

func (a *URLUnescape) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		return url.QueryUnescape(text)
	})
}
