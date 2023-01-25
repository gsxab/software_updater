package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
)

type ReadAttr struct {
	Default
	DefaultFactory[ReadAttr, *ReadAttr]
	Attribute string `json:"attribute"`
}

func (r *ReadAttr) Path() Path {
	return Path{"browser", "reader", "read_attr"}
}

func (a *ReadAttr) InElmNum() int {
	return 1
}

func (a *ReadAttr) InStrNum() int {
	return Any
}

func (a *ReadAttr) OutElmNum() int {
	return 1
}

func (a *ReadAttr) OutStrNum() int {
	return 1
}

func (a *ReadAttr) Do(ctx context.Context, _ selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	text, err := input.Elements[0].GetAttribute(a.Attribute)
	if err != nil {
		logs.Error(ctx, "selenium element get_attr failed", err, "attr", a.Attribute)
		return
	}
	output = StringToArgs(text, input)
	return
}

func (a *ReadAttr) ToDTO() *DTO {
	return &DTO{
		Input:  []string{"node"},
		Output: []string{"attribute_text"},
		Values: map[string]string{"attribute": a.Attribute},
	}
}
