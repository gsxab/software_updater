package hook

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/util/error_util"
)

type Variables struct {
	ActionPtr *action.Action
	Input     *action.Args
	Output    *action.Args
	ResultPtr *action.Result
	ErrorPtr  *error
}

type Hook struct {
	F    func(ctx context.Context, vars *Variables, id string, errs error_util.Collector)
	Name string
}

type RegisterInfo struct {
	Action   action.Path
	Hook     Hook
	Position *Position
	Event    Event
}
