package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type AppendFormat struct {
	Default
	DefaultFactory[AppendFormat, *AppendFormat]
	Format string `json:"format"`
}

func (a *AppendFormat) Path() Path {
	return Path{"string", "mutator", "append_format"}
}

func (a *AppendFormat) OutStrNum() int {
	return OneMore
}

func (a *AppendFormat) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	texts := make([]any, 0, len(input.Strings))
	for _, text := range input.Strings {
		texts = append(texts, text)
	}
	result := fmt.Sprintf(a.Format, texts...)
	output = AnotherStringToArgs(result, input)
	return
}

func (a *AppendFormat) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"formatted_text"},
		Values: map[string]string{"pattern": a.Format},
	}
}
