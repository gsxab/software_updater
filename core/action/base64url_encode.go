package action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type Base64URLEncode struct {
	StringMutator
	DefaultFactory[Base64URLEncode, *Base64URLEncode]
}

func (a *Base64URLEncode) Path() Path {
	return Path{"encoder", "rfc4648", "url_base64_encode"}
}

func (a *Base64URLEncode) Do(_ context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.Mutate(input, func(text string) string {
		return base64.URLEncoding.EncodeToString([]byte(text))
	})
}
