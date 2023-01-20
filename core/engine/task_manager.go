package engine

import (
	"context"
	"encoding/json"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/hook"
	"software_updater/core/job"
	"software_updater/core/util/error_util"
	"sync"
)

type JobManager struct {
}

func (t *JobManager) Resolve(ctx context.Context, actionStr string, actionManager *ActionManager, config *config.EngineConfig) (*job.Flow, error) {
	var storedFlow StoredFlow
	errs := error_util.NewCollector()
	errs.Collect(json.Unmarshal([]byte(actionStr), &storedFlow))

	flow := &job.Flow{
		Root: t.resolveBranch(ctx, storedFlow.Root, actionManager, config, errs),
	}

	return flow, errs.ToError()
}

func (t *JobManager) resolveBranch(ctx context.Context, storedBranch StoredBranch, actionManager *ActionManager,
	config *config.EngineConfig, errs *error_util.ErrorCollector) *job.Branch {
	jobs := make([]job.Job, 0, len(storedBranch.Actions))
	for _, storedAction := range storedBranch.Actions {
		a, hooks, err := actionManager.Action(storedAction.Path, storedAction.JSON)
		errs.Collect(err)
		jobs = append(jobs, t.NewJob(ctx, config, a, hooks))
	}

	next := make([]*job.Branch, 0, len(storedBranch.Next))
	for _, child := range storedBranch.Next {
		branch := t.resolveBranch(ctx, child, actionManager, config, errs)
		next = append(next, branch)
	}

	return &job.Branch{
		Jobs: jobs,
		Next: next,
	}
}

func (t *JobManager) NewJob(_ context.Context, config *config.EngineConfig, action action.Action, hooks []*hook.ActionHooks) job.Job {
	var item job.Job
	if config.DebugLog {
		item = &job.LoggedJob{}
	} else {
		item = &job.DefaultJob{}
	}
	item.SetAction(action, hooks)
	return item
}

func (t *JobManager) StaticCheckFlow(flow *job.Flow) error {
	var elmArgN, strArgN int
	errs := error_util.NewCollector()
	t.staticCheckBranch(flow.Root, elmArgN, strArgN, errs)
	return errs.ToError()
}

func (t *JobManager) staticCheckBranch(thread *job.Branch, elmArgN int, strArgN int, errs *error_util.ErrorCollector) {
	var err error
	for _, j := range thread.Jobs {
		elmArgN, err = action.StaticCheckInput(j.Action().InElmNum(), j.Action().OutElmNum(), elmArgN, j.Name())
		errs.Collect(err)
		strArgN, err = action.StaticCheckInput(j.Action().InStrNum(), j.Action().OutStrNum(), strArgN, j.Name())
		errs.Collect(err)
	}

	for _, child := range thread.Next {
		t.staticCheckBranch(child, elmArgN, strArgN, errs)
	}
}

func (t *JobManager) InitFlow(ctx context.Context, flow *job.Flow) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, 16)
	overChan := make(chan struct{})

	errs := error_util.NewCollector()
	go func() {
		err, ok := <-errChan
		for ok {
			errs.Collect(err)
			err, ok = <-errChan
		}
		close(overChan)
	}()

	t.initBranch(ctx, flow.Root, &error_util.ChannelCollector{errChan}, &wg)

	wg.Wait()
	close(errChan)
	<-overChan
	return errs.ToError()
}

func (t *JobManager) initBranch(ctx context.Context, branch *job.Branch, errs error_util.Collector, wg *sync.WaitGroup) {
	for _, j := range branch.Jobs {
		j.InitAction(ctx, errs, wg)
	}

	for _, child := range branch.Next {
		t.initBranch(ctx, child, errs, wg)
	}
}

func (t *JobManager) RunJobs(ctx context.Context, flow *job.Flow, driver selenium.WebDriver, cv *po.CurrentVersion) (*po.Version, error) {
	wg := sync.WaitGroup{}
	errChan := make(chan error, 16)
	overChan := make(chan struct{})

	errs := error_util.NewCollector()
	go func() {
		err, ok := <-errChan
		for ok {
			errs.Collect(err)
			err, ok = <-errChan
		}
		close(overChan)
	}()

	var args *action.Args
	v := cv.CurrentVersion

	t.runBranch(ctx, flow.Root, driver, args, v, &error_util.ChannelCollector{errChan}, &wg)

	wg.Wait()
	close(errChan)
	<-overChan
	return v, errs.ToError()
}

func (t *JobManager) runBranch(ctx context.Context, branch *job.Branch, driver selenium.WebDriver, args *action.Args,
	v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) {
	for _, j := range branch.Jobs {
		if args == nil {
			args = &action.Args{}
		}
		output, stop, _ := j.RunAction(ctx, driver, args, v, errs, wg)
		if stop {
			break
		}
		args = output
	}

	childCnt := len(branch.Next)
	if childCnt == 0 {
		return
	}
	wg.Add(childCnt - 1)
	for i := 1; i < childCnt; i++ {
		go func(branch2 *job.Branch) {
			t.runBranch(ctx, branch2, driver, args, v, errs, wg)
			wg.Done()
		}(branch.Next[i])
	}
	t.runBranch(ctx, branch.Next[0], driver, args, v, errs, wg)
}
