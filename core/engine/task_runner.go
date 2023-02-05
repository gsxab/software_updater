package engine

import (
	"context"
	cache "github.com/gsxab/go-generic_lru"
	cache2 "github.com/gsxab/go-generic_lru/lru_with_rw_lock"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/job"
	"software_updater/core/tools/web"
	"software_updater/core/util/error_util"
	"sync"
	"sync/atomic"
	"time"
)

type TaskID = int64

type Task struct {
	ID    TaskID
	State job.State
	Flow  *job.Flow
	CV    *po.CurrentVersion
	hpURL string
	// runtime
	Cancel func()
	Wg     *sync.WaitGroup
}

type TaskRunner struct {
	nextId  TaskID
	running *sync.Mutex                // locked if any task is running
	pending cache.Cache[TaskID, *Task] // with lock
	done    cache.Cache[TaskID, *Task] // with lock. less used than p so always lock in d -> p order.
	start   sync.Once
	close   <-chan struct{}
	dur     time.Duration
	update  func(ctx context.Context, cv *po.CurrentVersion, v *po.Version)
}

func NewTaskRunner(start bool, interval int, update func(context.Context, *po.CurrentVersion, *po.Version)) *TaskRunner {
	result := &TaskRunner{
		nextId:  0,
		pending: cache2.New[TaskID, *Task](0),
		done:    cache2.New[TaskID, *Task](config.Current().Engine.DoneCache),
		start:   sync.Once{},
		close:   make(chan struct{}),
		dur:     time.Duration(interval) * time.Second,
		update:  update,
	}
	if start {
		go result.StartRunningRoutine()
	}
	return result
}

func (r *TaskRunner) EnqueueJob(flow *job.Flow, cv *po.CurrentVersion, homepageURL string) (TaskID, error) {
	id := atomic.AddInt64(&r.nextId, 1)
	task := &Task{
		ID:    id,
		State: job.Init,
		Flow:  flow,
		CV:    cv,
	}
	r.pending.Add(id, task)
	task.State = job.Pending
	return id, nil
}

func (r *TaskRunner) GetTaskState(id TaskID) (bool, job.State, error) {
	var task *Task
	var ok bool

	task, ok = r.done.Get(id)
	if ok {
		return true, task.State, nil
	}

	task, ok = r.pending.Get(id)
	if ok {
		return true, task.State, nil
	}

	return false, 0, nil
}

func (r *TaskRunner) RunTask(ctx context.Context, task *Task, driver selenium.WebDriver, cv *po.CurrentVersion) (*po.Version, error) {
	r.running.Lock()
	defer r.running.Unlock()

	task.Wg = &sync.WaitGroup{}
	errChan := make(chan error, 16)
	errStopChan := make(chan struct{})

	errs := error_util.NewCollector()
	go func() {
		err, ok := <-errChan
		for ok {
			errs.Collect(err)
			err, ok = <-errChan
		}
		close(errStopChan)
	}()

	args := &action.Args{
		Elements: []selenium.WebElement{},
		Strings:  []string{task.hpURL},
	}
	now := time.Now()
	v := &po.Version{
		LocalTime:  &now,
		Filename:   nil,
		Picture:    nil,
		Link:       nil,
		Digest:     nil,
		RemoteDate: nil,
	}
	if cv.Version != nil {
		v.Name = cv.Version.Name
		v.Version = cv.Version.Version
		v.Previous = &cv.Version.ID
	}

	ctx, task.Cancel = context.WithCancel(ctx)
	task.State = job.Processing
	// multiple goroutines for var task, no write since here
	r.runBranch(ctx, task.Flow.Root, driver, args, v, &error_util.ChannelCollector{Channel: errChan}, task)
	task.Wg.Wait()
	// wait for goroutines for var task
	close(errChan)
	<-errStopChan

	err := errs.ToError()
	if task.State == job.Processing { // may be set to cancel
		if err != nil {
			task.State = job.Failure
		} else {
			task.State = job.Success
		}
	}
	if task.State != job.Success {
		return nil, err
	}

	return v, err
}

func (r *TaskRunner) runBranch(ctx context.Context, branch *job.Branch, driver selenium.WebDriver, args *action.Args,
	v *po.Version, errs error_util.Collector, task *Task) {
	for _, j := range branch.Jobs {
		if args == nil {
			args = &action.Args{}
		}
		output, finishBranch, stopFlow, _ := j.RunAction(ctx, driver, args, v, errs, task.Wg)
		if stopFlow {
			task.Cancel()
			task.State = job.Cancelled
			return
		}
		if finishBranch {
			return
		}
		args = output
	}

	childCnt := len(branch.Next)
	if childCnt == 0 {
		return
	}
	task.Wg.Add(childCnt - 1)
	for i := 1; i < childCnt; i++ {
		go func(branch2 *job.Branch) {
			r.runBranch(ctx, branch2, driver, args, v, errs, task)
			task.Wg.Done()
		}(branch.Next[i])
	}
	r.runBranch(ctx, branch.Next[0], driver, args, v, errs, task)
}

func (r *TaskRunner) StartRunningRoutine() {
	r.start.Do(func() {
		go func() {
			t := time.NewTimer(r.dur)
			defer func() {
				if !t.Stop() {
					<-t.C
				}
			}()
			ctx := context.Background()
		loop:
			for {
				select {
				case <-r.close:
					break loop
				case <-t.C:
					_, task, ok := r.pending.GetOldest()
					if !ok {
						t.Reset(r.dur)
						continue
					}

					version, err := r.RunTask(ctx, task, web.Driver(), task.CV)

					if version != nil {
						r.update(ctx, task.CV, version)
					}
					if err != nil {
						task.State = job.Failure
					}

					// lock both, order: d -> p
					r.done.ApplyRW(func(c cache.Cache[TaskID, *Task]) {
						r.pending.ApplyRW(func(c2 cache.Cache[TaskID, *Task]) {
							c2.Remove(task.ID)

							c.Add(task.ID, task)
							if task.State == job.Processing {
								task.State = job.Success
							}
						})
					})

					t.Reset(time.Second)
				}
			}
		}()
	})
}
