package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"strings"
	"sync"
)

type ReduceJoin struct {
	base.StringMutator
	base.DefaultFactory[ReduceJoin, *ReduceJoin]
	Sep  string `json:"sep"`
	Skip []int  `json:"skip,omitempty"`
}

func (a *ReduceJoin) Path() action.Path {
	return action.Path{"string", "mutator", "reduce_join"}
}

func (a *ReduceJoin) Icon() string {
	return "mdi:mdi-text-box-plus"
}

func (a *ReduceJoin) OutStrNum() int {
	return 1 + len(a.Skip)
}

func (a *ReduceJoin) Do(_ context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	results := make([]string, 1, len(input.Strings)+1)
	skipIndex := 0
	texts := make([]string, 1, len(input.Strings)+1)
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

	result := strings.Join(texts, a.Sep)
	results[0] = result
	output = action.StringsToArgs(results, input)
	return
}

func (a *ReduceJoin) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text"},
			Output: []string{"formatted_text"},
		},
		Values: map[string]string{"sep": a.Sep, "skip": util.ToJSON(a.Skip)},
	}
}
