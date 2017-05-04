package fibonacci

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBaseCases(t *testing.T) {
	result := CalculateFibonacci(1, 1)
	assert.Equal(t, big.NewInt(1), result)

	result = CalculateFibonacci(2, 1)
	assert.Equal(t, big.NewInt(1), result)
}

func TestIndex3(t *testing.T) {
	result := CalculateFibonacci(3, 1)
	assert.Equal(t, big.NewInt(2), result)
}

func TestIndex10(t *testing.T) {
	result := CalculateFibonacci(10, 1)
	assert.Equal(t, big.NewInt(55), result)
}

func TestIndex10Parallel(t *testing.T) {
	result := CalculateFibonacci(10, 4)
	assert.Equal(t, big.NewInt(55), result)
}

func TestIndex10MassivelyParallel(t *testing.T) {
	result := CalculateFibonacci(10, 100)
	assert.Equal(t, big.NewInt(55), result)
}
