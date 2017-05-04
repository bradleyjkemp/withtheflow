package concurrent_test

import (
	"github.com/bradleyjkemp/withtheflow/concurrent"
	"github.com/bradleyjkemp/withtheflow/examples/fibonacci"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestConcurrentFibonacci1Thread(t *testing.T) {
	suite.Run(t, fibonacci.CreateFibonacciTestSuite(concurrent.NewRunner(1)))
}

func TestConcurrentFibonacci2Threads(t *testing.T) {
	suite.Run(t, fibonacci.CreateFibonacciTestSuite(concurrent.NewRunner(2)))
}

func TestConcurrentFibonacci100Threads(t *testing.T) {
	suite.Run(t, fibonacci.CreateFibonacciTestSuite(concurrent.NewRunner(100)))
}
