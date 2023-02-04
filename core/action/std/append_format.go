package std

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"sync"
)

type AppendFormat struct {
	base.Default
	base.DefaultFactory[AppendFormat, *AppendFormat]
	Format string `json:"format"`
}

func (a *AppendFormat) Path() action.Path {
	return action.Path{"string", "mutator", "append_format"}
}

func (a *AppendFormat) Icon() string {
	return "text-box-plus-outline"
}

func (a *AppendFormat) OutStrNum() int {
	return action.OneMore
}

func (a *AppendFormat) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, wg *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	texts := make([]any, 0, len(input.Strings))
	for _, text := range input.Strings {
		texts = append(texts, text)
	}
	result := fmt.Sprintf(a.Format, texts...)
	output = action.AnotherStringToArgs(result, input)
	return
}

func (a *AppendFormat) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"formatted_text"},
		},
		Values: map[string]string{"pattern": a.Format},
	}
}
