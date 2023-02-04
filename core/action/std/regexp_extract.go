package std

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"regexp"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type RegexpExtract struct {
	base.StringMutator
	base.DefaultFactory[RegexpExtract, *RegexpExtract]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (a *RegexpExtract) Path() action.Path {
	return action.Path{"string", "mutator", "regexp_extract"}
}

func (a *RegexpExtract) Icon() string {
	return "regex"
}

func (a *RegexpExtract) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpExtract) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.MutateWithErr(ctx, input, func(text string) (string, error) {
		matched, result := util.MatchExtract(a.matcher, a.FullMatch, text)
		if !matched {
			return result, fmt.Errorf("matching failed, pattern: %s, text: %s", a.Pattern, text)
		}
		return result, nil
	})
}

func (a *RegexpExtract) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"text..."},
			Output: []string{"extracted_text..."},
		},
		Values: map[string]string{"pattern": a.Pattern, "skip": util.ToJSON(a.Skip)},
	}
}
