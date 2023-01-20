package action

import (
	"context"
	"encoding/json"
	"sync"
)

type Default struct {
}

func (d *Default) InElmNum() int {
	return Any
}

func (d *Default) InStrNum() int {
	return Any
}

func (d *Default) OutElmNum() int {
	return Same
}

func (d *Default) OutStrNum() int {
	return Same
}

func (d *Default) Init(context.Context, *sync.WaitGroup) error {
	return nil
}

type DefaultFactory[T any, PT interface {
	Action
	*T
}] struct{}

func (r *DefaultFactory[T, PT]) NewAction(args string) (Action, error) {
	ret := PT(new(T))
	err := json.Unmarshal([]byte(args), ret)
	return ret, err
}
