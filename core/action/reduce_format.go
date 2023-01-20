package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type ReduceFormat struct {
	StringMutator
	DefaultFactory[ReduceFormat, *ReduceFormat]
	Format string `json:"format"`
	Skip   []int  `json:"skip,omitempty"`
}

func (a *ReduceFormat) Path() Path {
	return Path{"string", " mutator", " reduce_format"}
}

func (a *ReduceFormat) OutStrNum() int {
	return 1 + len(a.Skip)
}

func (a *ReduceFormat) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
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
	output = StringsToArgs(results, input)
	return
}

func (a *ReduceFormat) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text"},
		Output: []string{"formatted_text"},
		Values: map[string]string{"pattern": a.Format, "skip": util.ToJSON(a.Skip)},
	}
}
