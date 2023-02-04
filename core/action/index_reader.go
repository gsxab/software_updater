package action

import (
	"context"
	"fmt"
	"software_updater/core/logs"
)

type IndexReader struct {
	Index int `json:"index"`
}

func (a *IndexReader) Read(ctx context.Context, input *Args, callback func(text string)) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	callback(text)
	output = input
	return
}

func (a *IndexReader) ReadDirectly(ctx context.Context, input *Args) (text string, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text = input.Strings[a.Index]
	return
}

func (a *IndexReader) ReadWithErr(ctx context.Context, input *Args, callback func(text string) error) (output *Args, exit Result, err error) {
	if len(input.Strings) <= a.Index {
		err = fmt.Errorf("array index out of bound, len: %d, index: %d", len(input.Strings), a.Index)
		logs.Error(ctx, "string slice indexing failed", err, "strings", input.Strings, "index", a.Index)
		return
	}
	text := input.Strings[a.Index]
	err = callback(text)
	output = input
	return
}
