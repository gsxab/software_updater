package engine

import (
	"context"
	"encoding/base64"
	cache "github.com/gsxab/go-generic_lru"
	cache2 "github.com/gsxab/go-generic_lru/lru_with_rw_lock"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/job"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util/error_util"
	"sync"
	"sync/atomic"
	"time"
)

type TaskID = int64

type Task struct {
	ID    TaskID
	Name  string
	State job.State
	Flow  *job.Flow
	CV    *po.CurrentVersion
	HpURL string
	// runtime
	Cancel context.CancelFunc
	Wg     *sync.WaitGroup
}

type TaskRunner struct {
	nextId  TaskID
	running *sync.Mutex                // locked if any task is running
	pending cache.Cache[TaskID, *Task] // with lock
	done    cache.Cache[TaskID, *Task] // with lock. less used than p so always lock in d -> p order.
	start   sync.Once
	dur     time.Duration
	update  func(ctx context.Context, cv *po.CurrentVersion, v *po.Version)
	cancel  context.CancelFunc
}

func NewTaskRunner(ctx context.Context, start bool, interval int, update func(context.Context, *po.CurrentVersion, *po.Version)) *TaskRunner {
	result := &TaskRunner{
		nextId:  0,
		running: &sync.Mutex{},
		pending: cache2.New[TaskID, *Task](0),
		done:    cache2.New[TaskID, *Task](config.Current().Engine.DoneCache),
		start:   sync.Once{},
		dur:     time.Duration(interval) * time.Second,
		update:  update,
	}
	if start {
		ctx, result.cancel = context.WithCancel(ctx)
		result.StartRunningRoutine(ctx)
	}
	return result
}

func (r *TaskRunner) EnqueueJob(ctx context.Context, flow *job.Flow, name string, cv *po.CurrentVersion, homepageURL string) (TaskID, error) {
	id := atomic.AddInt64(&r.nextId, 1)
	task := &Task{
		ID:    id,
		Name:  name,
		State: job.Init,
		Flow:  flow,
		CV:    cv,
		HpURL: homepageURL,
	}
	logs.InfoM(ctx, "enqueueing job", "id", id, "name", name)
	r.pending.Add(id, task)
	logs.InfoM(ctx, "enqueued job", "id", id, "name", name)
	task.State = job.Pending
	return id, nil
}

func (r *TaskRunner) Stop(ctx context.Context) {
	r.cancel()
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
	logs.InfoM(ctx, "starting task", "id", task.ID, "name", task.Name)
	defer logs.InfoM(ctx, "finished task", "id", task.ID)

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
		Strings:  []string{task.HpURL},
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
	if cv != nil && cv.Version != nil {
		v.Name = cv.Version.Name
		v.Version = cv.Version.Version
		v.Previous = &cv.Version.ID
	}

	ctx, task.Cancel = context.WithCancel(ctx)
	task.State = job.Processing

	var screenshot []byte
	{
		width, _ := driver.ExecuteScript("return document.body.scrollWidth;", nil)
		height, _ := driver.ExecuteScript("return document.body.scrollHeight;", nil)
		_ = driver.ResizeWindow("", int(width.(float64)), int(height.(float64)))
		screenshot, _ = driver.Screenshot()
	}

	// there may be multiple goroutines for vars, no write since here
	update := r.runBranch(ctx, task.Flow.Root, driver, args, v, &error_util.ChannelCollector{Channel: errChan}, task)
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

	if err != nil {
		return nil, err
	}
	if !update {
		return nil, nil
	}

	if v.Picture == nil {
		filename := base64.URLEncoding.EncodeToString([]byte(v.Name)) + "@" + time.Now().Format("2006-01-02") + ".png"
		err = os.WriteFile(path.Join(config.Current().Files.ScreenshotDir, filename), screenshot, os.FileMode(0o755))
		if err != nil {
			logs.Error(ctx, "write file failed", err)
			return nil, err
		}
		v.Picture = &filename
	}
	return v, nil
}

func (r *TaskRunner) runBranch(ctx context.Context, branch *job.Branch, driver selenium.WebDriver, args *action.Args,
	v *po.Version, errs error_util.Collector, task *Task) (update bool) {
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
	update = args.Update

	childCnt := len(branch.Next)
	if childCnt == 0 {
		return
	}
	task.Wg.Add(childCnt - 1)
	for i := 1; i < childCnt; i++ {
		go func(branch2 *job.Branch) {
			if r.runBranch(ctx, branch2, driver, args, v, errs, task) {
				update = true
			}
			task.Wg.Done()
		}(branch.Next[i])
	}
	if r.runBranch(ctx, branch.Next[0], driver, args, v, errs, task) {
		update = true
	}
	return
}

func (r *TaskRunner) StartRunningRoutine(ctx context.Context) {
	r.start.Do(func() {
		go func() {
			t := time.NewTimer(r.dur)
		loop:
			for {
				select {
				case <-ctx.Done():
					if !t.Stop() {
						<-t.C
					}
					break loop
				case <-t.C:
					func() {
						plan := r.dur
						defer func() {
							if err := recover(); err != nil {
								logs.ErrorM(context.Background(), "recovered failure", "err", err)
							}

							t.Reset(plan)
						}()

						_, task, ok := r.pending.GetOldest()
						if !ok {
							logs.InfoM(ctx, "no task is found pending now")
							return
						}

						version, err := r.RunTask(ctx, task, web.Driver(), task.CV)

						if version != nil {
							r.update(ctx, task.CV, version)
						}
						if err != nil {
							task.State = job.Failure
						}

						// lock both, order: d -> p
						r.done.ApplyRW(func(d cache.Cache[TaskID, *Task]) {
							r.pending.ApplyRW(func(p cache.Cache[TaskID, *Task]) {
								p.Remove(task.ID)

								d.Add(task.ID, task)
								if task.State == job.Processing {
									task.State = job.Success
								}
							})
						})

						plan = time.Second
					}()
				}
			}
		}()
	})
}
