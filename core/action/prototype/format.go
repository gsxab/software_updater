package prototype

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type Format struct {
	StringMutator
	DefaultFactory[Format, *Format]
	Format string `json:"format"`
}

func (a *Format) Path() action.Path {
	return action.Path{"string", "mutator", "format"}
}

func (a *Format) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Mutate(input, func(text string) string {
		return fmt.Sprintf(a.Format, text)
	})
}

func (a *Format) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"text"},
		Output: []string{"formatted_text"},
		Values: map[string]string{"pattern": a.Format, "skip": util.ToJSON(a.Skip)},
	}
}
