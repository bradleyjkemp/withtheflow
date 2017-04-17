package main

import (
	"fmt"
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/examples/fibonacci_iterative/proto/fibonacci"
	"github.com/bradleyjkemp/withtheflow/serial"
	"github.com/golang/protobuf/proto"
)

const (
	FibStep   = "fibStep"
	InitStep  = "init"
	PrintStep = "printStep"
)

func createFlowStep(index int64) withtheflow.FlowHandler {
	return func(data []byte, flowHandle withtheflow.FlowHandle) ([]byte, error) {
		fibData := &fibonacci.Data{}
		if err := proto.Unmarshal(data, fibData); err != nil {
			return nil, err
		}

		if fibData.Index == index {
			resultData := &fibonacci.Result{
				N:     fibData.N,
				Index: fibData.Index,
			}
			result, _ := proto.Marshal(resultData)

			flowHandle.NewFlow(withtheflow.FlowCall{PrintStep, result})

			return result, nil
		}

		fibData.N, fibData.NMinus1 = fibData.N+fibData.NMinus1, fibData.N
		fibData.Index++

		nextStepData, _ := proto.Marshal(fibData)
		flowHandle.NewFlow(withtheflow.FlowCall{FibStep, nextStepData})

		return nil, nil
	}
}

func initial(_ []byte, flowHandle withtheflow.FlowHandle) ([]byte, error) {
	fibData := &fibonacci.Data{
		NMinus1: 0,
		N:       1,
		Index:   1,
	}

	args, _ := proto.Marshal(fibData)

	flowHandle.NewFlow(withtheflow.FlowCall{FibStep, args})

	return nil, nil
}

func printResult(result []byte, _ withtheflow.FlowHandle) ([]byte, error) {
	resultData := &fibonacci.Result{}
	err := proto.Unmarshal(result, resultData)
	if err != nil {
		return nil, err
	}

	fmt.Printf("The %dth fibonacci number is %d\n", resultData.Index, resultData.N)

	return nil, nil
}

func main() {
	flowHandlers := map[string](withtheflow.FlowHandler){
		FibStep:   createFlowStep(10),
		InitStep:  initial,
		PrintStep: printResult,
	}

	workflow := serial.NewWorkflow(flowHandlers, nil)
	workflow.NewFlow(withtheflow.FlowCall{InitStep, nil})

	err := workflow.Run()

	if err != nil {
		fmt.Println(err)
	}

	graph, err := workflow.GenerateDotGraph()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(graph)
}
