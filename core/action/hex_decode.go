package action

import (
	"context"
	"encoding/hex"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type HexDecode struct {
	StringMutator
	DefaultFactory[HexDecode, *HexDecode]
}

func (a *HexDecode) Path() Path {
	return Path{"string", " decoder", " hex_decode"}
}

func (a *HexDecode) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(input, func(text string) (string, error) {
		bytes, err := hex.DecodeString(text)
		return string(bytes), err
	})
}
