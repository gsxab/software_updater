package std

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"sync"
)

type Base64URLDecode struct {
	prototype.StringMutator
	prototype.DefaultFactory[Base64URLDecode, *Base64URLDecode]
}

func (a *Base64URLDecode) Path() action.Path {
	return action.Path{"decoder", "rfc4648", "url_base64_decode"}
}

func (a *Base64URLDecode) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		bytes, err := base64.URLEncoding.DecodeString(text)
		return string(bytes), err
	})
}
