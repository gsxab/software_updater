package base

import (
	"context"
	"encoding/json"
	"software_updater/core/action"
	"sync"
)

type Default struct {
}

func (d *Default) Icon() string {
	return "ray-vertex"
}

func (d *Default) InElmNum() int {
	return action.Any
}

func (d *Default) InStrNum() int {
	return action.Any
}

func (d *Default) OutElmNum() int {
	return action.Same
}

func (d *Default) OutStrNum() int {
	return action.Same
}

func (d *Default) Init(context.Context, *sync.WaitGroup) error {
	return nil
}

type DefaultFactory[T any, PT interface {
	action.Action
	*T
}] struct{}

func (r *DefaultFactory[T, PT]) NewAction(args string) (action.Action, error) {
	ret := PT(new(T))
	err := json.Unmarshal([]byte(args), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *DefaultFactory[T, PT]) ToProtoDTO() *action.ProtoDTO {
	t := PT(new(T))
	return t.ToDTO().ProtoDTO
}
