package concurrent

import (
	"github.com/bradleyjkemp/withtheflow"
	"sync"
	"sync/atomic"
)

type executionToken struct{}

var flowIdCounter int64

type workflow struct {
	handlers       map[string]withtheflow.FlowHandler
	results        map[int64]*flowResult
	resultsMutex   sync.Mutex
	executionSlots chan executionToken
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
	return &workflow{
		handlers:       handlers,
		results:        make(map[int64]*flowResult),
		executionSlots: make(chan executionToken, concurrency),
	}
}

func (w *workflow) Run(funcname string, args interface{}) interface{} {
	runtime := &workflowRuntime{w}
	id := runtime.AddFlow(funcname, args)

	return runtime.getResult(id.(int64))
}

func (w *workflowRuntime) createFlow() int64 {
	result := &flowResult{}
	result.waiter.Add(1)
	flowId := atomic.AddInt64(&flowIdCounter, 1)

	w.resultsMutex.Lock()
	w.results[flowId] = result
	w.resultsMutex.Unlock()

	return flowId
}

func (w *workflowRuntime) getResult(flowId int64) interface{} {
	w.resultsMutex.Lock()
	result := w.results[flowId]
	w.resultsMutex.Unlock()
	result.waiter.Wait()

	if result.deferredResult != 0 {
		return w.getResult(result.deferredResult)
	}

	return result.result
}

func (w *workflowRuntime) setResult(flowId int64, flowResult interface{}) {
	w.resultsMutex.Lock()
	result := w.results[flowId]
	result.result = flowResult
	result.waiter.Done()
	w.resultsMutex.Unlock()
}

func (w *workflowRuntime) setDeferredResult(flowId int64, deferredId int64) {
	w.resultsMutex.Lock()
	result := w.results[flowId]
	result.deferredResult = deferredId
	result.waiter.Done()
	w.resultsMutex.Unlock()
}

func (w *workflowRuntime) AddFlow(funcname string, args interface{}, dependentIds ...withtheflow.FlowId) withtheflow.FlowId {
	flowId := w.createFlow()

	go func() {
		var results []interface{}
		for _, flowId := range dependentIds {
			results = append(results, w.getResult(flowId.(int64)))
		}

		w.executionSlots <- executionToken{}
		result := w.handlers[funcname](args, w, results)

		if r, ok := result.(deferredResult); ok {
			w.setDeferredResult(flowId, r.deferredId)
		} else {
			w.setResult(flowId, result)
		}
		<-w.executionSlots
	}()

	return flowId
}

func (w *workflowRuntime) DeferredResult(deferredId withtheflow.FlowId) interface{} {
	return deferredResult{deferredId.(int64)}
}
