package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"sync"
)

type ReadAttr struct {
	Default
	DefaultFactory[ReadAttr, *ReadAttr]
	Attribute string `json:"attribute"`
}

func (r *ReadAttr) Path() Path {
	return Path{"node", "reader", "read_attr"}
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

func (a *ReadAttr) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	text, err := input.Elements[0].GetAttribute(a.Attribute)
	if err != nil {
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
