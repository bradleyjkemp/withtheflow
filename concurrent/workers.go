package concurrent

import (
	"context"
)

func (w *workflowRuntime) spawnWorker(ctx context.Context) {
	go func() {
		var job int64
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.idQueue.GetTaskChan():
				job = w.idQueue.GetTask()
				w.executeFlow(job)
			}
		}
	}()
}

func (w *workflowRuntime) executeFlow(jobId int64) {
	w.mutex.Lock()
	job := w.tasks[jobId]
	w.mutex.Unlock()

	if len(job.dependentIds) > 0 {
		w.executeDependentFlow(jobId, job)
		return
	}

	w.executionSlots <- struct{}{}
	result := w.handlers[job.funcname](job.args, w, nil)
	w.setResult(jobId, result)
	<-w.executionSlots
}

func (w *workflowRuntime) executeDependentFlow(jobId int64, job *flowTask) {
	go func() {
		var results []interface{}
		for _, dependentId := range job.dependentIds {
			results = append(results, w.getResult(dependentId.(int64)))
		}

		w.executionSlots <- struct{}{}
		result := w.handlers[job.funcname](job.args, w, results)
		w.setResult(jobId, result)
		<-w.executionSlots
	}()
}
