package main

import (
	"fmt"
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/examples/fibonacci_recursive/proto/fibonacci"
	"github.com/bradleyjkemp/withtheflow/serial"
	"github.com/golang/protobuf/proto"
)

const (
	FibStep      = "FibStep"
	AddReducer   = "AddReducer"
	PrintReducer = "PrintReducer"
	InitStep     = "InitStep"
)

// FibStep(takes {n:int}) returns {n:int}
// calls BlockedResult([FibStep(n-1), FibStep(n-2)], AddStep)

func fibStep(args []byte, flowHandle withtheflow.FlowHandle) ([]byte, error) {
	argsData := &fibonacci.FibStepArg{}
	if err := proto.Unmarshal(args, argsData); err != nil {
		return nil, err
	}

	if argsData.Index == 1 || argsData.Index == 2 {
		resultData := &fibonacci.Result{
			N: 1,
		}

		return proto.Marshal(resultData)
	}

	argsMinus1, _ := proto.Marshal(&fibonacci.FibStepArg{argsData.Index - 1})
	nMinus1, _ := flowHandle.NewFlow(withtheflow.FlowCall{FibStep, argsMinus1})

	argsMinus2, _ := proto.Marshal(&fibonacci.FibStepArg{argsData.Index - 2})
	nMinus2, _ := flowHandle.NewFlow(withtheflow.FlowCall{FibStep, argsMinus2})

	flowHandle.NewBlockedResult(withtheflow.FlowReducerCall{AddReducer, nil, []withtheflow.FlowId{nMinus1, nMinus2}})

	return nil, nil
}

func addReducer(results []withtheflow.FlowResult, _ []byte) ([]byte, error) {
	total := &fibonacci.Result{}

	for _, result := range results {
		resultData := &fibonacci.Result{}
		proto.Unmarshal(result.Result, resultData)

		total.N += resultData.N
	}

	return proto.Marshal(total)
}

func printReducer(results []withtheflow.FlowResult, args []byte) ([]byte, error) {
	resultData := &fibonacci.Result{}
	proto.Unmarshal(results[0].Result, resultData)

	argsData := &fibonacci.FibStepArg{}
	proto.Unmarshal(args, argsData)

	fmt.Printf("The %dth fibonacci number is %d\n", argsData.Index, resultData.N)

	return nil, nil
}

func initStep(args []byte, flowHandle withtheflow.FlowHandle) ([]byte, error) {
	argsData := &fibonacci.FibStepArg{}
	argsData.Index = 20

	args, _ = proto.Marshal(argsData)

	calc, _ := flowHandle.NewFlow(withtheflow.FlowCall{FibStep, args})

	flowHandle.NewBlockedResult(withtheflow.FlowReducerCall{PrintReducer, args, []withtheflow.FlowId{calc}})

	return nil, nil
}

func main() {
	flowHandlers := map[string](withtheflow.FlowHandler){
		FibStep:  fibStep,
		InitStep: initStep,
	}

	flowReducers := map[string](withtheflow.FlowReducer){
		AddReducer:   addReducer,
		PrintReducer: printReducer,
	}

	workflow := serial.NewWorkflow(flowHandlers, flowReducers)

	workflow.NewFlow(withtheflow.FlowCall{InitStep, nil})

	err := workflow.Run()

	if err != nil {
		fmt.Println(err)
	}

	// graph, err := workflow.GenerateDotGraph()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(graph)
}
