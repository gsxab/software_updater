package hook

import (
	"fmt"
	"software_updater/core/util/slice_util"
)

type Event string

const (
	BeforeInitEvent Event = "before_init"
	AfterInitEvent  Event = "after_init"
	BeforeRunEvent  Event = "before"
	AfterRunEvent   Event = "after"
)

type ActionHooks struct {
	BeforeInit []Hook
	AfterInit  []Hook
	BeforeRun  []Hook
	AfterRun   []Hook
}

func (h *ActionHooks) HookPtr(event Event) *[]Hook {
	switch event {
	case BeforeInitEvent:
		return &h.BeforeInit
	case AfterInitEvent:
		return &h.AfterInit
	case BeforeRunEvent:
		return &h.BeforeRun
	case AfterRunEvent:
		return &h.AfterRun
	default:
		panic("unexpected event")
	}
}

func (h *ActionHooks) Get(event Event) []Hook {
	return *h.HookPtr(event)
}

func (h *ActionHooks) PutAt(event Event, hook Hook, pos *Position) error {
	ptr := h.HookPtr(event)

	switch pos.Cmd {
	case FirstCmd:
		*ptr = slice_util.Prepend(*ptr, hook)
	case LastCmd:
		*ptr = append(*ptr, hook)
	case PrevCmd:
		in, idx := slice_util.LinearSearchWithPtr(*ptr, func(x *Hook) bool { return x.Name == pos.Ref })
		if !in {
			return fmt.Errorf("registered not found")
		}
		*ptr = slice_util.Insert(*ptr, idx, hook)
	case NextCmd:
		in, idx := slice_util.LinearSearchWithPtr(*ptr, func(x *Hook) bool { return x.Name == pos.Ref })
		if !in {
			return fmt.Errorf("registered not found")
		}
		*ptr = slice_util.Insert(*ptr, idx+1, hook)
	default:
		return fmt.Errorf("unknown position command")
	}
	return nil
}
