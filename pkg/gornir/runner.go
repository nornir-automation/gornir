package gornir

import (
	"sync"
)

// JobParams is an interface that needs to be implemented by the parameters object passed to the jobs
type JobParams interface{}

// Task is a function supported by the Runner facility
type Task interface {
	Run(Context, *sync.WaitGroup, chan *JobResult)
}

// TBD
type Runner interface {
	Run(Context, Task, chan *JobResult) error
	Close() error
	Wait() error
}

// JobResult is an interface that needs to be implemented by the job result object
type JobResult struct {
	ctx        Context
	err        error
	changed    bool
	data       interface{}
	subResults []*JobResult
}

func NewJobResult(ctx Context) *JobResult {
	return &JobResult{ctx: ctx}
}

func (r *JobResult) Context() Context {
	return r.ctx
}

func (r *JobResult) SetContext(ctx Context) {
	r.ctx = ctx
}

func (r *JobResult) Err() error {
	return r.err
}

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

func (r *JobResult) SetErr(err error) {
	r.err = err
}

func (r *JobResult) Changed() bool {
	return r.changed
}

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

func (r *JobResult) SetChanged(changed bool) {
	r.changed = changed
}

func (r *JobResult) Data() interface{} {
	return r.data
}

func (r *JobResult) SetData(data interface{}) {
	r.data = data
}

func (r *JobResult) SubResults() []*JobResult {
	return r.subResults
}

func (r *JobResult) AddSubResult(result *JobResult) {
	r.subResults = append(r.subResults, result)
}
