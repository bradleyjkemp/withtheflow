package fibonacci

import (
	"github.com/bradleyjkemp/withtheflow"
	"math/big"
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

func CalculateFibonacci(runner withtheflow.WorkflowRunner, index int) *big.Int {
	return runner.SetFlowHandlers(map[string]withtheflow.FlowHandler{
		FIB_STEP: fibStep,
		ADD_STEP: addStep,
	}).Run(FIB_STEP, index).(*big.Int)
}
