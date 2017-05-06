package mandelbrot

import (
	"bytes"
	"fmt"
	"github.com/bradleyjkemp/withtheflow/concurrent"
	"github.com/bradleyjkemp/withtheflow/mocks"
	"github.com/stretchr/testify/assert"
	"math/cmplx"
	"testing"
)

func generateRow(y float64) []bool {
	var row []bool
	for x := -2.0; x <= 0.5; x += 0.0315 {
		row = append(row, cmplx.Abs(mandelbrot(complex(x, y))) < 2)
	}
	return row
}

// Adapted from rosettacode.org/wiki/Mandelbrot_set#Go
func referenceGenerator() string {
	var output bytes.Buffer

	for y := 1.0; y >= -1.0; y -= 0.05 {
		for x := -2.0; x <= 0.5; x += 0.0315 {
			if cmplx.Abs(mandelbrot(complex(x, y))) < 2 {
				output.WriteString("*")
			} else {
				output.WriteString(" ")
			}
		}
		output.WriteString("\n")
	}

	return output.String()
}

func convertToImage(in [][]bool) string {
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

func TestRow(t *testing.T) {
	actualRow := calculateRow(row{
		X_MIN,
		0,
		(X_MAX - X_MIN) / float64(80),
		80,
	}, &mocks.Runtime{}, nil)

	expectedRow := generateRow(0)

	assert.Equal(t, expectedRow, actualRow)
}

func TestImage(t *testing.T) {
	boolImage := GenerateMandelbrot(concurrent.NewRunner(1), 80, 41)

	actual := convertToImage(boolImage)
	fmt.Println(len(actual))
	expected := referenceGenerator()
	fmt.Println(len(expected))
	assert.Equal(t, expected, actual)
}
