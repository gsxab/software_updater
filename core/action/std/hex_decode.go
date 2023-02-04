package std

import (
	"context"
	"encoding/hex"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"
)

type HexDecode struct {
	base.StringMutator
	base.DefaultFactory[HexDecode, *HexDecode]
}

func (a *HexDecode) Path() action.Path {
	return action.Path{"string", "decoder", "hex_decode"}
}

func (a *HexDecode) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		bytes, err := hex.DecodeString(text)
		return string(bytes), err
	})
}
