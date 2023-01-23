package action

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"regexp"
	"software_updater/core/db/po"
	"software_updater/core/util"
	"sync"
)

type RegexpExtract struct {
	StringMutator
	DefaultFactory[RegexpExtract, *RegexpExtract]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (a *RegexpExtract) Path() Path {
	return Path{"string", "mutator", "regexp_extract"}
}

func (a *RegexpExtract) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpExtract) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	return a.MutateWithErr(input, func(text string) (string, error) {
		matched, result := util.MatchExtract(a.matcher, a.FullMatch, text)
		if !matched {
			return result, fmt.Errorf("matching failed, pattern: %s, text: %s", a.Pattern, text)
		}
		return result, nil
	})
}

func (a *RegexpExtract) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"text..."},
		Output: []string{"extracted_text..."},
		Values: map[string]string{"pattern": a.Pattern, "skip": util.ToJSON(a.Skip)},
	}
}
