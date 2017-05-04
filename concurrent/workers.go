package concurrent

import (
	"context"
)

func (w *workflowRuntime) spawnWorker(ctx context.Context) {
	go func() {
		var job *flowTask
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.getStackSlot:
				// fmt.Println("Got task slot")
				job = w.popTask()
				w.executeFlow(job)
			}
		}
	}()
}

func (w *workflowRuntime) executeFlow(job *flowTask) {
	if len(job.dependentIds) > 0 {
		w.executeDependentFlow(job)
		return
	}

	w.executionSlots <- struct{}{}
	result := w.handlers[job.funcname](job.args, w, nil)
	w.setResult(job.flowId, result)
	<-w.executionSlots
}

func (w *workflowRuntime) executeDependentFlow(job *flowTask) {
	go func() {
		var results []interface{}
		for _, dependentId := range job.dependentIds {
			results = append(results, w.getResult(dependentId.(int64)))
		}

		w.executionSlots <- struct{}{}
		result := w.handlers[job.funcname](job.args, w, results)
		w.setResult(job.flowId, result)
		<-w.executionSlots
	}()
}
