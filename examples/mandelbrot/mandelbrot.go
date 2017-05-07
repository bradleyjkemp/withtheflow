package mandelbrot

import (
	"github.com/bradleyjkemp/withtheflow"
	"math/cmplx"
)

const (
	APPENDER      = "appender"
	CALCULATE_ROW = "calculateRow"
	GENERATOR     = "generateMandelbrot"
	X_MIN         = -2.0
	X_MAX         = -2.0 + 80*0.0315
	Y_MIN         = -1
	Y_MAX         = 1
	ITERATIONS    = 50
)

type row struct {
	// (x,y) of top leftmost corner of the row
	x float64
	y float64

	// resolution of the image
	stepSize float64

	// dimensions of the row in pixels
	width int
}

type dimensions struct {
	x int
	y int
}

func mandelbrot(a complex128) (z complex128) {
	for i := 0; i < 50; i++ {
		z = z*z + a
	}
	// fmt.Printf("Testing %s and got %d\n", a, cmplx.Abs(z))
	return
}

func calculateRow(args interface{}, _ withtheflow.Runtime, _ []interface{}) interface{} {
	rowInfo := args.(row)

	var row []bool
	for xSteps := 0; xSteps < rowInfo.width; xSteps++ {
		x := float64(rowInfo.x) + (float64(xSteps) * rowInfo.stepSize)
		abs := cmplx.Abs(mandelbrot(complex(x, rowInfo.y)))
		row = append(row, abs < 2)
	}

	return row
}

// appends rows together to form the image
func rowAppender(args interface{}, _ withtheflow.Runtime, rows []interface{}) interface{} {
	var image [][]bool
	for _, row := range rows {
		singleRow := row.([]bool)
		image = append(image, singleRow)
	}

	return image
}

func generateMandelbrot(args interface{}, runtime withtheflow.Runtime, _ []interface{}) interface{} {
	dimension := args.(dimensions)
	var rows []withtheflow.FlowId
	xPitch := (X_MAX - X_MIN) / float64(dimension.x)
	yPitch := (Y_MAX - Y_MIN) / float64(dimension.y)

	for ySteps := 0; ySteps < dimension.y; ySteps++ {
		y := Y_MIN + (float64(ySteps) * yPitch)
		rowInfo := row{
			x:        X_MIN,
			y:        y,
			stepSize: xPitch,
			width:    dimension.x,
		}
		rows = append(rows, runtime.AddFlow(CALCULATE_ROW, rowInfo))
	}

	return runtime.DeferredResult(runtime.AddFlow(APPENDER, nil, rows...))
}

func GenerateMandelbrot(runner withtheflow.WorkflowRunner, x, y int) [][]bool {
	return runner.SetFlowHandlers(map[string]withtheflow.FlowHandler{
		APPENDER:      rowAppender,
		CALCULATE_ROW: calculateRow,
		GENERATOR:     generateMandelbrot,
	}).Run(GENERATOR, dimensions{x, y}).([][]bool)
}
