package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type Format struct {
	StringMutator
	DefaultFactory[Format, *Format]
	Format string `json:"format"`
}

func (a *Format) Path() Path {
	return Path{"string", "mutator", "format"}
}

func (a *Format) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.Mutate(input, func(text string) string {
		return fmt.Sprintf(a.Format, text)
	})
}

func (a *Format) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"formatted_text"},
		Values: map[string]string{"pattern": a.Format, "skip": util.ToJSON(a.Skip)},
	}
}
