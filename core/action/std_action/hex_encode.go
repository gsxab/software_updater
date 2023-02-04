package std_action

import (
	"context"
	"encoding/hex"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"sync"
)

type HexEncode struct {
	action.StringMutator
	action.DefaultFactory[HexEncode, *HexEncode]
}

func (a *HexEncode) Path() action.Path {
	return action.Path{"encoder", "rfc4648", "hex_encode"}
}

func (a *HexEncode) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Mutate(input, func(text string) string {
		return hex.EncodeToString([]byte(text))
	})
}
