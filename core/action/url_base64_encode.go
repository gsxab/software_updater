package action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type URLBase64Encode struct {
	StringMutator
	DefaultFactory[URLBase64Encode, *URLBase64Encode]
}

func (a *URLBase64Encode) Path() Path {
	return Path{"encoder", " rfc4648", " url_base64_encode"}
}

func (a *URLBase64Encode) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.Mutate(input, func(text string) string {
		return base64.URLEncoding.EncodeToString([]byte(text))
	})
}
