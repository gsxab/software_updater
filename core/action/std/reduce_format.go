package std

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/prototype"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type ReduceFormat struct {
	prototype.StringMutator
	prototype.DefaultFactory[ReduceFormat, *ReduceFormat]
	Format string `json:"format"`
	Skip   []int  `json:"skip,omitempty"`
}

func (a *ReduceFormat) Path() action.Path {
	return action.Path{"string", "mutator", "reduce_format"}
}

func (a *ReduceFormat) OutStrNum() int {
	return 1 + len(a.Skip)
}

func (a *ReduceFormat) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	results := make([]string, 1, len(input.Strings)+1)
	skipIndex := 0
	texts := make([]any, 1, len(input.Strings)+1)
	for index, text := range input.Strings {
		// strings skipped will be pushed into results[1:]
		if skipIndex < len(a.Skip) && index == a.Skip[skipIndex] {
			skipIndex++
			results = append(results, text)
			continue
		}
		// the remaining will be formatted
		texts = append(texts, text)
	}

	result := fmt.Sprintf(a.Format, texts...)
	results[0] = result
	output = action.StringsToArgs(results, input)
	return
}

func (a *ReduceFormat) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"text"},
		Output: []string{"formatted_text"},
		Values: map[string]string{"pattern": a.Format, "skip": util.ToJSON(a.Skip)},
	}
}
