package main

import (
	"fmt"
	"github.com/bradleyjkemp/withtheflow/examples/fibonacci_iterative/workflow"
)

func main() {
	var result int64
	var index int64 = 10

	workflow := workflow.CreateWorkflow(index, &result)
	err := workflow.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("The %dth fibonacci number is %d\n", index, result)
}
