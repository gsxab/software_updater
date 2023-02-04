package std_action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type ReadAttr struct {
	action.Default
	action.DefaultFactory[ReadAttr, *ReadAttr]
	Attribute string `json:"attribute"`
}

func (r *ReadAttr) Path() action.Path {
	return action.Path{"browser", "reader", "read_attr"}
}

func (a *ReadAttr) InElmNum() int {
	return 1
}

func (a *ReadAttr) InStrNum() int {
	return action.Any
}

func (a *ReadAttr) OutElmNum() int {
	return 1
}

func (a *ReadAttr) OutStrNum() int {
	return 1
}

func (a *ReadAttr) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	text, err := input.Elements[0].GetAttribute(a.Attribute)
	if err != nil {
		logs.Error(ctx, "selenium element get_attr failed", err, "attr", a.Attribute)
		return
	}
	output = action.StringToArgs(text, input)
	return
}

func (a *ReadAttr) ToDTO() *action.DTO {
	return &action.DTO{
		Input:  []string{"node"},
		Output: []string{"attribute_text"},
		Values: map[string]string{"attribute": a.Attribute},
	}
}
