/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Copyright (c) 2023. gsxab.
 *
 * This file is part of Software Update Watcher, a.k.a. Zhixin Robot.
 *
 * Software Update Watcher is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
 *
 * Software Update Watcher is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with Software Update Watcher. If not, see <https://www.gnu.org/licenses/>.
 */

package engine

import (
	"container/list"
	"context"
	"errors"
	"reflect"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/flow"
	"software_updater/core/tools/web"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gsxab/go-error_util/errcollect"
	cache "github.com/gsxab/go-generic_lru"
	cache_impl "github.com/gsxab/go-generic_lru/lru_with_rw_lock"
	"github.com/gsxab/go-logs"
	"github.com/tebeka/selenium"
)

type TaskID = int64

type Task struct {
	ID    TaskID
	Name  string
	State flow.State
	Flow  *flow.Flow
	CV    *po.CurrentVersion
	HpURL string
	// runtime
	Cancel context.CancelFunc
	Wg     *sync.WaitGroup
}

func (t *Task) MetaDTO() *flow.TaskMetaDTO {
	return &flow.TaskMetaDTO{
		ID:      t.ID,
		Name:    t.Name,
		State:   t.State.Int(),
		Version: t.CV.Version.Version,
	}
}

func (t *Task) DTO() *flow.TaskDTO {
	return &flow.TaskDTO{
		DTO:         t.Flow.ToDTO(),
		TaskMetaDTO: t.MetaDTO(),
	}
}

type TaskValuer interface {
	Value() *Task
}

type TaskRunner struct {
	nextId  TaskID
	running *sync.Mutex                // locked if any task is running
	pending cache.Cache[TaskID, *Task] // with lock
	done    cache.Cache[TaskID, *Task] // with lock. less used than p so always lock in d -> p order.
	nameMap map[string]TaskID
	start   sync.Once
	dur     time.Duration
	update  func(ctx context.Context, cv *po.CurrentVersion, v *po.Version)
	cancel  context.CancelFunc
}

func NewTaskRunner(ctx context.Context, start bool, interval int, update func(context.Context, *po.CurrentVersion, *po.Version)) *TaskRunner {
	nameMap := make(map[string]TaskID)
	result := &TaskRunner{
		nextId:  0,
		running: &sync.Mutex{},
		pending: cache_impl.New[TaskID, *Task](0),
		done: cache_impl.NewWithEvicted[TaskID, *Task](config.Current().Engine.DoneCache, func(id TaskID, task *Task) {
			delete(nameMap, task.Name)
		}),
		nameMap: nameMap,
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

func (r *TaskRunner) EnqueueJob(ctx context.Context, fl *flow.Flow, name string, cv *po.CurrentVersion, homepageURL string) (TaskID, error) {
	if taskID, ok := r.nameMap[name]; ok {
		if _, found := r.pending.Get(taskID); found {
			logs.WarnM(ctx, "enqueueing a job duplicated with a pending one, skipping")
			return taskID, nil
		}
	}

	id := atomic.AddInt64(&r.nextId, 1)
	task := &Task{
		ID:    id,
		Name:  name,
		State: flow.Init,
		Flow:  fl,
		CV:    cv,
		HpURL: homepageURL,
	}
	logs.InfoM(ctx, "enqueueing flow", "id", id, "name", name)
	r.pending.Add(id, task)
	logs.InfoM(ctx, "enqueued flow", "id", id, "name", name)
	task.State = flow.Pending
	r.nameMap[name] = id
	return id, nil
}

func (r *TaskRunner) Stop(ctx context.Context) {
	r.cancel()
}

func (r *TaskRunner) GetTask(id TaskID) (bool, *Task, error) {
	var task *Task
	var ok bool

	task, ok = r.pending.Get(id)
	if ok {
		return true, task, nil
	}

	task, ok = r.done.Get(id)
	if ok {
		return true, task, nil
	}

	return false, nil, nil
}

func (r *TaskRunner) GetAllTasks(ctx context.Context) (tasks []*Task, err error) {
	r.done.ApplyRO(func(d cache.Cache[TaskID, *Task]) {
		r.pending.ApplyRO(func(p cache.Cache[TaskID, *Task]) {
			var listAny interface{}

			listAny, err = d.Container()
			if err != nil {
				return
			}
			taskList, ok := listAny.(*list.List)
			if !ok {
				logs.ErrorM(ctx, "d.container is not *list.List", "type", reflect.TypeOf(listAny))
				err = errors.New("intenal error")
				return
			}
			for e := taskList.Back(); e != nil; e = e.Prev() {
				task, ok := e.Value.(TaskValuer)
				if !ok {
					logs.ErrorM(ctx, "one of the entries in p.container is not evaluated to Task", "type", reflect.TypeOf(e.Value))
					continue
				}
				tasks = append(tasks, task.Value())
			}

			listAny, err = p.Container()
			if err != nil {
				return
			}
			taskList, ok = listAny.(*list.List)
			if !ok {
				logs.ErrorM(ctx, "p.container is not *list.List", "type", reflect.TypeOf(listAny))
				err = errors.New("intenal error")
				return
			}
			for e := taskList.Back(); e != nil; e = e.Prev() {
				task, ok := e.Value.(TaskValuer)
				if !ok {
					logs.ErrorM(ctx, "one of the entries in p.container is not evaluated to Task", "type", reflect.TypeOf(e.Value))
					continue
				}
				tasks = append(tasks, task.Value())
			}
		})
	})

	return tasks, nil
}

func (r *TaskRunner) GetTaskIDMap() (map[string]TaskID, error) {
	return r.nameMap, nil
}

func (r *TaskRunner) RunTask(ctx context.Context, task *Task, driver selenium.WebDriver, cv *po.CurrentVersion) (*po.Version, error) {
	r.running.Lock()
	defer r.running.Unlock()
	logs.InfoM(ctx, "starting task", "id", task.ID, "name", task.Name)
	defer logs.InfoM(ctx, "finished task", "id", task.ID)

	task.Wg = &sync.WaitGroup{}
	errChan := make(chan error, 16)
	errStopChan := make(chan struct{})

	errs := errcollect.New()
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
	task.State = flow.Processing

	// there may be multiple goroutines for vars, no write since here
	update := r.runBranch(ctx, task.Flow.Root, driver, args, v, errcollect.NewFromChannel(errChan), task)
	task.Wg.Wait()
	// wait for goroutines for var task
	close(errChan)
	<-errStopChan

	err := errs.ToError()
	if task.State == flow.Processing { // may be set to cancel
		if err != nil {
			task.State = flow.Failure
		} else {
			task.State = flow.Success
		}
	}
	if task.State != flow.Success {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if !update {
		return nil, nil
	}

	if v.Picture == nil {
		filename, err := web.TakeScreenshot(ctx, driver, v.Name)
		if err != nil {
			logs.Warn(ctx, "ignoring error in taking a screenshot", err)
		} else {
			v.Picture = &filename
		}
	}
	return v, nil
}

func (r *TaskRunner) runBranch(ctx context.Context, branch *flow.Branch, driver selenium.WebDriver, args *action.Args,
	v *po.Version, errs errcollect.Collector, task *Task) (update bool) {
	for _, j := range branch.Steps {
		if args == nil {
			args = &action.Args{}
		}
		output, finishBranch, stopFlow, _ := j.RunAction(ctx, driver, args, v, errs, task.Wg)
		if stopFlow {
			task.Cancel()
			task.State = flow.Cancelled
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
		go func(branch2 *flow.Branch) {
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
							if msg := recover(); msg != nil {
								logs.ErrorM(ctx, "recovered failure", "msg", msg)
							}

							t.Reset(plan)
						}()

						_, task, ok := r.pending.GetOldest()
						if !ok {
							logs.InfoM(ctx, "no task is found pending now")
							return
						}
						err := web.Driver().Get("about:blank")
						if err != nil {
							logs.Warn(ctx, "go to blank page failure", err)
							err = nil
						}

						version, err := r.RunTask(ctx, task, web.Driver(), task.CV)

						if version != nil {
							r.update(ctx, task.CV, version)
						}
						if err != nil {
							task.State = flow.Failure
						}

						// lock both, order: d -> p
						r.done.ApplyRW(func(d cache.Cache[TaskID, *Task]) {
							r.pending.ApplyRW(func(p cache.Cache[TaskID, *Task]) {
								p.Remove(task.ID)

								d.Add(task.ID, task)
								if task.State == flow.Processing {
									task.State = flow.Success
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
