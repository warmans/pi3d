package main

import (
	"fmt"
	"github.com/deadsy/sdfx/sdf"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	PiFile           = "10000.txt"
	LineLength       = 32
	NumLines         = 32
	CubeSize         = 5.0
	BaseHeight       = 5.0

	// exaggerate the height of the cubes by multiplying them by this number.
	HeightMultiplier = 2

	// Oversize the blocks slightly to allow them to be merged together more cleanly.
	OverlapFactor = 1.2 //e.g. 1.2 = 120%
)

func main() {

	grid := make([][]sdf.SDF3, 0)

	pi := mustGetPi()
	for i := 0; i < LineLength*NumLines; i++ {

		if len(grid) > NumLines {
			break
		}
		if len(grid) == 0 || len(grid[len(grid)-1]) == LineLength {
			grid = append(grid, []sdf.SDF3{})
		}

		var height = 0.0
		if len(pi) > i {
			height = float64(pi[i])
		}
		height = height * HeightMultiplier

		// create the cube
		c := sdf.Box3D(
			sdf.V3{
				X: CubeSize * OverlapFactor,
				Y: CubeSize * OverlapFactor,
				// make the cubes extend half way into the base to make the union cleaner.
				Z: height + (BaseHeight / 2),
			},
			0,
		)
		c = sdf.Transform3D(
			c,
			sdf.Translate3d(
				sdf.V3{
					X: CubeSize*float64(len(grid[len(grid)-1])) + (CubeSize / 2),
					Y: CubeSize*float64(len(grid)-1) + (CubeSize / 2),
					Z: 0 - (BaseHeight - height/2) + BaseHeight,
				},
			),
		)
		grid[len(grid)-1] = append(grid[len(grid)-1], c)
	}

	//create a base
	base := sdf.Box3D(
		sdf.V3{
			// base also needs to be oversized
			X: (CubeSize * OverlapFactor) * LineLength,
			Y: (CubeSize * OverlapFactor) * NumLines,
			Z: BaseHeight,
		},
		0,
	)
	base = sdf.Transform3D(
		base,
		sdf.Translate3d(sdf.V3{
			X: CubeSize * float64(LineLength) / 2,
			Y: CubeSize * float64(NumLines) / 2,
			Z: 0 - (BaseHeight / 2),
		}),
	)

	//join everything together
	for _, line := range grid {
		base = sdf.Union3D(base, sdf.Union3D(line...))
	}

	sdf.RenderSTL(base, 250, "pi.stl")
}

func mustGetPi() []int {
	f, err := os.Open(PiFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	ints := []int{}

	str := fmt.Sprintf("%s", b)
	beforeDecimalPoint := true
	for _, char := range str {
		if string(char) == "." {
			beforeDecimalPoint = false
			continue
		}
		intVal, err := strconv.Atoi(string(char))
		if err != nil {
			panic(err)
		}
		if beforeDecimalPoint {
			ints = append(ints, intVal*10)
		} else {
			ints = append(ints, intVal)
		}
	}
	return ints
}
