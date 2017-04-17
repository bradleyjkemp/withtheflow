package workflow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test10(t *testing.T) {
	var result int64
	var index int64 = 10

	workflow := CreateWorkflow(index, &result)
	err := workflow.Run()
	assert.NoError(t, err)

	assert.Equal(t, int64(55), result)
}

func Test20(t *testing.T) {
	var result int64
	var index int64 = 20

	workflow := CreateWorkflow(index, &result)
	err := workflow.Run()
	assert.NoError(t, err)

	assert.Equal(t, int64(6765), result)
}
