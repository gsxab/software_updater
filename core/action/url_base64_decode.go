package action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type URLBase64Decode struct {
	StringMutator
	DefaultFactory[URLBase64Decode, *URLBase64Decode]
}

func (a *URLBase64Decode) Path() Path {
	return Path{"decoder", "rfc4648", "url_base64_decode"}
}

func (a *URLBase64Decode) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		bytes, err := base64.URLEncoding.DecodeString(text)
		return string(bytes), err
	})
}
