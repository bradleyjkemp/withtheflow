package concurrent

import (
	"fmt"
	"github.com/bradleyjkemp/withtheflow"
)

var stackSize int

type flowTask struct {
	flowId       int64
	funcname     string
	args         interface{}
	dependentIds []withtheflow.FlowId
}

func setupInfiniteChannel(in chan struct{}, out chan struct{}) {
	go func() {
		for {
			out <- <-in
		}
	}()
}

func (w *workflowRuntime) pushTask(task *flowTask) {
	// 	fmt.Println("Pushing")
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.addStackSlot <- struct{}{}
	w.jobStack = append(w.jobStack, task)
	stackSize++
	fmt.Printf("Pushed %d\n", stackSize)
}

func (w *workflowRuntime) popTask() *flowTask {
	w.mutex.Lock()
	stackSize--
	fmt.Printf("Popped %d\n", stackSize)
	defer w.mutex.Unlock()

	stackSize := len(w.jobStack)
	task := w.jobStack[stackSize-1]
	w.jobStack = w.jobStack[:stackSize-1]
	return task
}
