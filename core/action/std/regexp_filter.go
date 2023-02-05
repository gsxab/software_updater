package std

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"regexp"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"sync"
)

type RegexpFilter struct {
	base.DefaultFactory[RegexpFilter, *RegexpFilter]
	Pattern string `json:"pattern"`
	matcher *regexp.Regexp
}

func (r *RegexpFilter) Path() action.Path {
	return action.Path{"selector", "filter", "regexp_filter"}
}

func (a *RegexpFilter) Icon() string {
	return "filter-outline"
}

func (a *RegexpFilter) InElmNum() int {
	return action.Any
}

func (a *RegexpFilter) InStrNum() int {
	return action.Any
}

func (a *RegexpFilter) OutElmNum() int {
	return 1
}

func (a *RegexpFilter) OutStrNum() int {
	return action.Same
}

func (a *RegexpFilter) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpFilter) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements := input.Elements
	var text string
	for _, element := range elements {
		text, err = element.Text()
		if err != nil {
			logs.Error(ctx, "selenium element get_text failed", err)
			return
		}
		if a.matcher != nil {
			matched := util.Match(a.matcher, true, text)
			if !matched {
				continue
			}
			output = action.ElementToArgs(element, input)
		}
	}
	err = fmt.Errorf("find matching element failed, matcher: %s, elements: %v", a.Pattern, elements)
	logs.Error(ctx, "element matching failed", err)
	return
}

func (a *RegexpFilter) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"nodes..."},
			Output: []string{"node"},
		},
		Values: map[string]string{"pattern": a.Pattern},
	}
}