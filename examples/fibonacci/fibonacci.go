package fibonacci

import (
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/concurrent"
	"math/big"
	"time"
)

const (
	FIB_STEP = "fibStep"
	ADD_STEP = "addStep"
)

func addStep(_ interface{}, runtime withtheflow.Runtime, subResults []interface{}) interface{} {
	total := &big.Int{}
	for _, arg := range subResults {
		total.Add(total, arg.(*big.Int))
	}

	// sleep so that most of the "work" is done in the flow handler rather than the flow runner
	time.Sleep(100 * time.Millisecond)
	return total
}

func fibStep(arg interface{}, runtime withtheflow.Runtime, _ []interface{}) interface{} {
	index := arg.(int)
	if index == 1 || index == 2 {
		return big.NewInt(1)
	} else {
		nMinus1 := runtime.AddFlow(FIB_STEP, index-1)
		nMinus2 := runtime.AddFlow(FIB_STEP, index-2)

		sum := runtime.AddFlow(ADD_STEP, nil, nMinus1, nMinus2)

		return runtime.DeferredResult(sum)
	}
}

func CalculateFibonacci(index int, concurrency int) *big.Int {
	return concurrent.NewWorkflow(map[string]withtheflow.FlowHandler{
		FIB_STEP: fibStep,
		ADD_STEP: addStep,
	}, concurrency).Run(FIB_STEP, index).(*big.Int)
}
