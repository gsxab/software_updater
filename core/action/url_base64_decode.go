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

func (a *URLBase64Decode) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(input, func(text string) (string, error) {
		bytes, err := base64.URLEncoding.DecodeString(text)
		return string(bytes), err
	})
}
