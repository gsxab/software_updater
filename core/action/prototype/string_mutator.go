package prototype

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/core/util/error_util"
)

type StringMutator struct {
	Default
	Skip []int `json:"skip,omitempty"`
}

func (a *StringMutator) Mutate(input *action.Args, mutate func(text string) string) (output *action.Args, exit action.Result, err error) {
	results := make([]string, 0, len(input.Strings))
	skipIndex := 0
	for index, text := range input.Strings {
		if skipIndex < len(a.Skip) && index == a.Skip[skipIndex] {
			skipIndex++
			results = append(results, text)
			continue
		}

		result := mutate(text)
		results = append(results, result)
	}
	output = action.StringsToArgs(results, input)
	return
}

func (a *StringMutator) MutateWithErr(ctx context.Context, input *action.Args, mutate func(text string) (string, error)) (output *action.Args, exit action.Result, err error) {
	errs := error_util.NewCollector()
	output, exit, err = a.Mutate(input, func(text string) string {
		result, err := mutate(text)
		errs.CollectWithLog(err, func(err error) {
			logs.Error(ctx, "string mutating failed", err, "input", text)
		})
		return result
	})
	errs.Collect(err)
	return output, exit, errs.ToError()
}

func (a *StringMutator) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"text..."},
		Output: []string{"formatted_text..."},
		Values: map[string]string{"skip": util.ToJSON(a.Skip)},
	}
}
