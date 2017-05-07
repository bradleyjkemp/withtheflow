package mandelbrot

import (
	"bytes"
	"github.com/bradleyjkemp/withtheflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/cmplx"
)

// Adapted from rosettacode.org/wiki/Mandelbrot_set#Go
func referenceGenerator() [][]bool {
	var img [][]bool

	for y := 1.0; y >= -1.0; y -= 0.05 {
		var row []bool
		for x := -2.0; x <= 0.5; x += 0.0315 {
			row = append(row, cmplx.Abs(mandelbrot(complex(x, y))) < 2)
		}
		img = append(img, row)
	}

	return img
}

func convertToAscii(in [][]bool) string {
	var output bytes.Buffer
	for _, row := range in {
		for _, pixel := range row {
			if pixel {
				output.WriteString("*")
			} else {
				output.WriteString(" ")
			}
		}
		output.WriteString("\n")
	}
	return output.String()
}

type MandelbrotTestSuite struct {
	suite.Suite
	runner func() withtheflow.WorkflowRunner
}

func (s *MandelbrotTestSuite) TestImage() {
	boolActual := GenerateMandelbrot(s.runner(), 80, 40)
	boolExpected := referenceGenerator()
	assert.Equal(s.T(), boolExpected, boolActual)
}

func CreateTestSuite(runner func() withtheflow.WorkflowRunner) *MandelbrotTestSuite {
	testSuite := new(MandelbrotTestSuite)
	testSuite.runner = runner
	return testSuite
}
