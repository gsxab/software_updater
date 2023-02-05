package job

import (
	"fmt"
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/util/error_util"
	"sync"
)

type DefaultJob struct {
	name      string
	action    action.Action
	hooks     []*hook.ActionHooks
	state     State
	stateDesc string
}

func (j *DefaultJob) SetAction(action action.Action, hooks []*hook.ActionHooks) {
	j.action = action
	j.hooks = hooks
}

func (j *DefaultJob) Action() action.Action {
	return j.action
}

func (j *DefaultJob) InitAction(ctx context.Context, errs error_util.Collector, wg *sync.WaitGroup) {
	j.state = Init

	select {
	case <-ctx.Done():
		j.state = Cancelled
		return
	default:
	}
	// ready

	// hooks: before init
	for _, actionHooks := range j.hooks {
		hooks := actionHooks.BeforeInit
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{}, j.name, errs)
		}
	}

	// init action
	err := j.action.Init(ctx, wg)

	// hooks: after init
	for i := len(j.hooks) - 1; i >= 0; i++ {
		hooks := j.hooks[i].AfterInit
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ErrorPtr: &err,
			}, j.name, errs)
		}
	}

	// after init
	if err != nil {
		j.state = Failure
		errs.Collect(fmt.Errorf("action init error: %w", err))
		return
	}
	j.state = Ready
}

func (j *DefaultJob) RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) (output *action.Args, finishBranch bool, stopFlow bool, err error) {
	j.state = Processing

	select {
	case <-ctx.Done():
		j.state = Cancelled
		return nil, true, true, nil
	default:
	}
	// ready

	// hooks: before run
	for _, actionHooks := range j.hooks {
		hooks := actionHooks.BeforeRun
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
				Input:     args,
			}, j.name, errs)
		}
	}

	// run action
	output, result, err := j.action.Do(ctx, driver, args, v, wg)

	// hooks: after run
	for i := len(j.hooks) - 1; i >= 0; i++ {
		hooks := j.hooks[i].AfterRun
		for _, h := range hooks {
			h.F(ctx, &hook.Variables{
				ActionPtr: &j.action,
				Input:     args,
				Output:    output,
				ResultPtr: &result,
				ErrorPtr:  &err,
			}, j.name, errs)
		}
	}

	// after run
	if err != nil {
		j.state = Aborted
		errs.Collect(fmt.Errorf("action run error: %w", err))
		return output, true, true, err
	}
	switch result {
	case action.Cancelled:
		j.state = Cancelled
		finishBranch = true
	case action.Again:
		return j.RunAction(ctx, nil, args, v, errs, nil)
	case action.StopBranch:
		j.state = Success
		finishBranch = true
	case action.StopFlow:
		j.state = Success
		finishBranch = true
		stopFlow = true
	case action.Finished:
		j.state = Success
	case action.Skipped:
		j.state = Skipped
	default:
		errs.Collect(fmt.Errorf("invalid action state: %d", result))
	}
	return
}

func (j *DefaultJob) State() State {
	return j.state
}

func (j *DefaultJob) SetState(state State) {
	j.state = state
}

func (j *DefaultJob) SetStateDesc(s string) {
	j.stateDesc = s
}

func (j *DefaultJob) Name() string {
	return j.name
}

func (j *DefaultJob) SetName(name string) {
	j.name = name
}

func (j *DefaultJob) ToDTO() *DTO {
	return &DTO{
		State:     j.state.Int(),
		StateDesc: j.stateDesc,
	}
}
