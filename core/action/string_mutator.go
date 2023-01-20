package action

import (
	"software_updater/core/util"
	"software_updater/core/util/error_util"
)

type StringMutator struct {
	Default
	Skip []int `json:"skip,omitempty"`
}

func (a *StringMutator) Mutate(input *Args, mutate func(text string) string) (output *Args, exit Result, err error) {
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
	output = StringsToArgs(results, input)
	return
}

func (a *StringMutator) MutateWithErr(input *Args, mutate func(text string) (string, error)) (output *Args, exit Result, err error) {
	errs := error_util.NewCollector()
	output, exit, err = a.Mutate(input, func(text string) string {
		result, err := mutate(text)
		errs.Collect(err)
		return result
	})
	errs.Collect(err)
	return output, exit, errs.ToError()
}

func (a *StringMutator) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text..."},
		Output: []string{"formatted_text..."},
		Values: map[string]string{"skip": util.ToJSON(a.Skip)},
	}
}
