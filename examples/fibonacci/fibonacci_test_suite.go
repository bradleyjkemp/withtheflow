package fibonacci

import (
	"github.com/bradleyjkemp/withtheflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/big"
)

type FibonacciTestSuite struct {
	suite.Suite
	runner withtheflow.WorkflowRunner
}

func (s *FibonacciTestSuite) TestBaseCases() {
	result := CalculateFibonacci(s.runner, 1)
	assert.Equal(s.T(), big.NewInt(1), result)

	result = CalculateFibonacci(s.runner, 2)
	assert.Equal(s.T(), big.NewInt(1), result)
}

func (s *FibonacciTestSuite) TestIndex3() {
	result := CalculateFibonacci(s.runner, 3)
	assert.Equal(s.T(), big.NewInt(2), result)
}

func (s *FibonacciTestSuite) TestIndex10() {
	result := CalculateFibonacci(s.runner, 10)
	assert.Equal(s.T(), big.NewInt(55), result)
}

func CreateFibonacciTestSuite(runner withtheflow.WorkflowRunner) *FibonacciTestSuite {
	testSuite := new(FibonacciTestSuite)
	testSuite.runner = runner
	return testSuite
}
