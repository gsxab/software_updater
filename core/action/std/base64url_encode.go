package std

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"
)

type Base64URLEncode struct {
	base.StringMutator
	base.DefaultFactory[Base64URLEncode, *Base64URLEncode]
}

func (a *Base64URLEncode) Path() action.Path {
	return action.Path{"encoder", "rfc4648", "url_base64_encode"}
}

func (a *Base64URLEncode) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Mutate(input, func(text string) string {
		return base64.URLEncoding.EncodeToString([]byte(text))
	})
}
