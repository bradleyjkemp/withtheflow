package concurrent

import (
	"context"
	"github.com/bradleyjkemp/withtheflow"
	"sync"
	"sync/atomic"
)

const (
	STACK_BUFFER_SIZE = 10
)

var flowIdCounter int64

type workflow struct {
	handlers         map[string]withtheflow.FlowHandler
	results          map[int64]*flowResult
	mutex            sync.Mutex
	executionSlots   chan struct{}
	jobStack         []*flowTask
	addStackSlot     chan struct{}
	getStackSlot     chan struct{}
	workerCancelFunc context.CancelFunc
}

type flowResult struct {
	result         interface{}
	deferredResult int64
	// The waiter represents whether the result (or location of the result) of this flow is known
	// Only the flow with the associated id should modify this struct and once done() is
	// called no further changes should be made.
	waiter sync.WaitGroup
}

type workflowRuntime struct {
	*workflow
}

type deferredResult struct {
	deferredId int64
}

func NewWorkflow(handlers map[string]withtheflow.FlowHandler, concurrency int) withtheflow.WorkflowRunner {
	w := &workflow{
		handlers:       handlers,
		results:        make(map[int64]*flowResult),
		executionSlots: make(chan struct{}, concurrency),
		addStackSlot:   make(chan struct{}, STACK_BUFFER_SIZE),
		getStackSlot:   make(chan struct{}, STACK_BUFFER_SIZE),
	}
	setupInfiniteChannel(w.addStackSlot, w.getStackSlot)

	return w
}

func (w *workflow) Run(funcname string, args interface{}) interface{} {
	runtime := &workflowRuntime{w}
	ctx, cancel := context.WithCancel(context.Background())
	w.workerCancelFunc = cancel

	for i := 0; i < cap(w.executionSlots); i++ {
		runtime.spawnWorker(ctx)
	}

	id := runtime.AddFlow(funcname, args)

	return runtime.getResult(id.(int64))
}

func (w *workflow) Close() {
	w.workerCancelFunc()
}

func (w *workflowRuntime) createFlow() int64 {
	result := &flowResult{}
	result.waiter.Add(1)
	flowId := atomic.AddInt64(&flowIdCounter, 1)

	w.mutex.Lock()
	w.results[flowId] = result
	w.mutex.Unlock()

	return flowId
}

func (w *workflowRuntime) getResult(flowId int64) interface{} {
	w.mutex.Lock()
	result := w.results[flowId]
	w.mutex.Unlock()
	result.waiter.Wait()

	w.deleteResult(flowId)

	if result.deferredResult != 0 {
		return w.getResult(result.deferredResult)
	}

	return result.result
}

func (w *workflowRuntime) deleteResult(flowId int64) {
	w.mutex.Lock()
	// A result can only ever be read once so, now that we've read it,
	// delete it from the map to limit memory usage
	delete(w.results, flowId)
	w.mutex.Unlock()
}

func (w *workflowRuntime) setResult(flowId int64, flowResult interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	result := w.results[flowId]
	if r, ok := flowResult.(deferredResult); ok {
		result.deferredResult = r.deferredId
	} else {
		result.result = flowResult
	}
	result.waiter.Done()
}

func (w *workflowRuntime) AddFlow(funcname string, args interface{}, dependentIds ...withtheflow.FlowId) withtheflow.FlowId {
	flowId := w.createFlow()

	job := &flowTask{
		flowId:       flowId,
		funcname:     funcname,
		args:         args,
		dependentIds: dependentIds,
	}
	w.pushTask(job)
	return flowId
}

func (w *workflowRuntime) DeferredResult(deferredId withtheflow.FlowId) interface{} {
	return deferredResult{deferredId.(int64)}
}
