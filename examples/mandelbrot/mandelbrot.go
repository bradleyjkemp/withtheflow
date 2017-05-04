package mandelbrot

import (
	"github.com/bradleyjkemp/withtheflow"
	"math/cmplx"
)

const (
	APPENDER        = "appender"
	CALCULATE_CHUNK = "calculateChunk"
)

type row struct {
	// (x,y) of top leftmost corner of the row
	x float64
	y float64

	// resolution of the image
	stepSize float64

	// dimensions of the row in pixels
	width  int
	height int
}

func calculateRow(args interface{}, _ withtheflow.Runtime, _ []interface{}) interface{} {
	chunkInfo := args.(*row)

	var chunk [][]bool

	for ySteps := 0; ySteps < chunkInfo.height; ySteps++ {
		y := chunkInfo.y + (float64(ySteps) * chunkInfo.stepSize)
		var row []bool
		for xSteps := 0; xSteps < chunkInfo.width; xSteps++ {
			x := chunkInfo.x + (float64(xSteps) * chunkInfo.stepSize)

			a := complex(x, y)
			var z complex128
			for i := 0; i < 50; i++ {
				z = z*z + a
			}
			if cmplx.Abs(z) < 2 {
				row[xSteps] = true
			}
		}
		chunk = append(chunk, row)
	}

	return chunk
}

// appends rows together to form the image
func rowAppender(args interface{}, _ withtheflow.Runtime, rows []interface{}) interface{} {
	var image [][]bool
	for _, row := range rows {
		singleRows := row.([][]bool)
		image = append(image, singleRows...)
	}

	return image
}
