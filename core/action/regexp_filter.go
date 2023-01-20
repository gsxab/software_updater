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

type RegexpFilter struct {
	DefaultFactory[RegexpFilter, *RegexpFilter]
	Pattern string `json:"pattern"`
	matcher *regexp.Regexp
}

func (r *RegexpFilter) Path() Path {
	return Path{"selector", "filter", "regexp_filter"}
}

func (a *RegexpFilter) InElmNum() int {
	return Any
}

func (a *RegexpFilter) InStrNum() int {
	return Any
}

func (a *RegexpFilter) OutElmNum() int {
	return 1
}

func (a *RegexpFilter) OutStrNum() int {
	return Same
}

func (a *RegexpFilter) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpFilter) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	elements := input.Elements
	var text string
	for _, element := range elements {
		text, err = element.Text()
		if err != nil {
			return
		}
		if a.matcher != nil {
			matched := util.Match(a.matcher, true, text)
			if !matched {
				continue
			}
			output = ElementToArgs(element, input)
		}
	}
	err = fmt.Errorf("find matching element failed, matcher: %s, elements: %v", a.Pattern, elements)
	return
}

func (a *RegexpFilter) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"nodes..."},
		Output: []string{"node"},
		Values: map[string]string{"pattern": a.Pattern},
	}
}
