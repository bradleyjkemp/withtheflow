package concurrent

import (
	"github.com/bradleyjkemp/withtheflow"
)

var stackSize int

type flowTask struct {
	flowId       int64
	funcname     string
	args         interface{}
	dependentIds []withtheflow.FlowId
}

func setupInfiniteChannel(inChan chan struct{}, outChan chan struct{}) {
	var slots int
	in := inChan
	out := outChan
	go func() {
		for {
			if slots == 0 {
				out = nil
			} else {
				out = outChan
			}

			select {
			case <-in:
				slots++
			case out <- struct{}{}:
				slots--
			}
		}
	}()
}

func (w *workflowRuntime) pushTask(task *flowTask) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.addStackSlot <- struct{}{}
	w.jobStack = append(w.jobStack, task)
	stackSize++
}

func (w *workflowRuntime) popTask() *flowTask {
	w.mutex.Lock()
	stackSize--
	defer w.mutex.Unlock()

	stackSize := len(w.jobStack)
	task := w.jobStack[stackSize-1]
	w.jobStack = w.jobStack[:stackSize-1]

	return task
}
