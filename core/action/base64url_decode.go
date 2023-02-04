package action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type Base64URLDecode struct {
	StringMutator
	DefaultFactory[Base64URLDecode, *Base64URLDecode]
}

func (a *Base64URLDecode) Path() Path {
	return Path{"decoder", "rfc4648", "url_base64_decode"}
}

func (a *Base64URLDecode) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		bytes, err := base64.URLEncoding.DecodeString(text)
		return string(bytes), err
	})
}
