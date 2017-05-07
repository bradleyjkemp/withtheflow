package concurrent_test

import (
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/concurrent"
	"github.com/bradleyjkemp/withtheflow/examples/fibonacci"
	"github.com/bradleyjkemp/withtheflow/examples/mandelbrot"
	"github.com/stretchr/testify/suite"
	"testing"
)

var threadTests = []int{1, 2, 100}

func runnerFactory(concurrency int) func() withtheflow.WorkflowRunner {
	return func() withtheflow.WorkflowRunner { return concurrent.NewRunner(concurrency) }
}

func TestConcurrentFibonacci(t *testing.T) {
	for _, concurrency := range threadTests {
		suite.Run(t, fibonacci.CreateTestSuite(runnerFactory(concurrency)))
	}
}

func TestConcurrentMandelbrot(t *testing.T) {
	for _, concurrency := range threadTests {
		suite.Run(t, mandelbrot.CreateTestSuite(runnerFactory(concurrency)))
	}
}
