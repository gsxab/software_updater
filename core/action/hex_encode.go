package action

import (
	"context"
	"encoding/hex"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type HexEncode struct {
	StringMutator
	DefaultFactory[HexEncode, *HexEncode]
}

func (a *HexEncode) Path() Path {
	return Path{"encoder", " rfc4648", " hex_encode"}
}

func (a *HexEncode) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.Mutate(input, func(text string) string {
		return hex.EncodeToString([]byte(text))
	})
}
