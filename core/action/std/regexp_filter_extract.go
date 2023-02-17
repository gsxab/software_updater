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

type RegexpFilterExtract struct {
	base.DefaultFactory[RegexpFilterExtract, *RegexpFilterExtract]
	Pattern   string `json:"pattern"`
	FullMatch bool   `json:"full_match"`
	matcher   *regexp.Regexp
}

func (r *RegexpFilterExtract) Path() action.Path {
	return action.Path{"selector", "filter", "regexp_filter_extract"}
}

func (a *RegexpFilterExtract) Icon() string {
	return "mdi:mdi-regex"
}

func (a *RegexpFilterExtract) InElmNum() int {
	return action.Any
}

func (a *RegexpFilterExtract) InStrNum() int {
	return action.Any
}

func (a *RegexpFilterExtract) OutElmNum() int {
	return 1
}

func (a *RegexpFilterExtract) OutStrNum() int {
	return 1
}

func (a *RegexpFilterExtract) Init(context.Context, *sync.WaitGroup) (err error) {
	a.matcher, err = regexp.Compile(a.Pattern)
	return
}

func (a *RegexpFilterExtract) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	elements := input.Elements
	var text string
	for _, element := range elements {
		text, err = element.Text()
		if err != nil {
			logs.Error(ctx, "selenium element get_text failed", err)
			return
		}
		matched, results := util.MatchExtractMultiple(a.matcher, a.FullMatch, text)
		if matched {
			output = action.StringsToArgs(results, input)
			output.Elements = []selenium.WebElement{element}
			return
		}
	}
	err = fmt.Errorf("find matching element failed, matcher: %s, elements: %v", a.Pattern, elements)
	logs.Error(ctx, "element matching failed", err)
	return
}

func (a *RegexpFilterExtract) ToDTO() *action.DTO {
	return &action.DTO{
		ProtoDTO: &action.ProtoDTO{
			Input:  []string{"nodes..."},
			Output: []string{"text"},
		},
		Values: map[string]string{"pattern": a.Pattern},
	}
}
