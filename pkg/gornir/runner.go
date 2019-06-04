package gornir

import (
	"context"
	"sync"
)

type TaskParameters struct {
	title  string
	logger Logger
	host   *Host
}

func NewTaskParameters(title string, logger Logger) *TaskParameters {
	return &TaskParameters{
		title:  title,
		logger: logger,
	}
}

func (tp *TaskParameters) Title() string {
	return tp.title
}

func (tp *TaskParameters) Logger() Logger {
	return tp.logger
}

func (tp *TaskParameters) Host() *Host {
	return tp.host
}

func (tp *TaskParameters) ForHost(host *Host) *TaskParameters {
	return &TaskParameters{
		title:  tp.title,
		logger: tp.logger.WithField("host", host.Hostname),
		host:   host,
	}
}

// Task is the interface that task plugins need to implement.
// the task is responsible to indicate its completion
// by calling sync.WaitGroup.Done()
type Task interface {
	Run(context.Context, *sync.WaitGroup, *TaskParameters, chan *JobResult)
}

// Runner is the interface of a struct that can implement a strategy
// to run tasks over hosts
type Runner interface {
	Run(context.Context, Task, map[string]*Host, *TaskParameters, chan *JobResult) error // Run executes the task over the hosts
	Close() error                                                                        // Close closes and cleans all objects associated with the runner
	Wait() error                                                                         // Wait blocks until all the hosts are done executing the task
}

// JobResult is the result of running a task over a host.
type JobResult struct {
	ctx        context.Context
	tp         *TaskParameters
	err        error
	changed    bool
	data       interface{}
	subResults []*JobResult
}

// NewJobResult instantiates a new JobResult
func NewJobResult(ctx context.Context, tp *TaskParameters) *JobResult {
	return &JobResult{
		ctx: ctx,
		tp:  tp,
	}
}

// Context returns the context associated with the task
func (r *JobResult) Context() context.Context {
	return r.ctx
}

func (r *JobResult) TaskParameters() *TaskParameters {
	return r.tp
}

// Err returns the error the task set, otherwise nil
func (r *JobResult) Err() error {
	return r.err
}

// AnyErr will return either the error the task set or any error reported by
// any subtask
func (r *JobResult) AnyErr() error {
	if r.err != nil {
		return r.err
	}
	for _, s := range r.subResults {
		if s.err != nil {
			return s.err
		}
	}
	return nil
}

// SetErr stores the error  and also propagates it to the associated Host
func (r *JobResult) SetErr(err error) {
	r.err = err
	r.TaskParameters().Host().setErr(err)
}

// Changed will return whether the task changed something or not
func (r *JobResult) Changed() bool {
	return r.changed
}

// AnyChanged will return whether the task or any of its subtasks
// changed something or not
func (r *JobResult) AnyChanged() bool {
	if r.changed {
		return true
	}
	for _, s := range r.subResults {
		if s.changed {
			return true
		}
	}
	return false
}

// SetChanged stores whether the task changed something or not
func (r *JobResult) SetChanged(changed bool) {
	r.changed = changed
}

// Data retrieves arbitrary data stored in the object
func (r *JobResult) Data() interface{} {
	return r.data
}

// SetData let's you store arbitrary data in the object
func (r *JobResult) SetData(data interface{}) {
	r.data = data
}

// SubResults returns the result of subtasks
func (r *JobResult) SubResults() []*JobResult {
	return r.subResults
}

// AddSubResult allows you to store the result of running a subtask
func (r *JobResult) AddSubResult(result *JobResult) {
	r.subResults = append(r.subResults, result)
}
