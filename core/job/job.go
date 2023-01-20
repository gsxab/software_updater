package job

import (
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/util/error_util"
	"sync"
)

type State int

const (
	Init State = iota + 1
	Pending
	Processing
	Success
	Failure
	Cancelled
	Skipped
	Aborted
)

func (s State) Int() int {
	return int(s)
}

type Job interface {
	SetAction(action action.Action, hooks []*hook.ActionHooks)
	Action() action.Action
	InitAction(ctx context.Context, errs error_util.Collector, wg *sync.WaitGroup)
	RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) (*action.Args, bool, error)
	State() State
	SetState(State)
	SetStateDesc(string)
	Name() string
	SetName(string)
	ToDTO() *DTO
}

type Branch struct {
	Jobs []Job
	Next []*Branch
}

type Flow struct {
	Root *Branch
}

type DebugInfo struct {
	Err    error
	Input  *action.Args
	Output *action.Args
}

type DTO struct {
	ActionDTO *action.DTO `json:"action"`
	State     int         `json:"state"`
	StateDesc string      `json:"state_desc"`
	DebugInfo *DebugInfo  `json:"debug_info"`
}
