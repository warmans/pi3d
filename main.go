package main

import (
	"flag"
	"fmt"
	"github.com/deadsy/sdfx/sdf"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	LineLength = 32
	NumLines   = 32
	CubeSize   = 5.0
	BaseHeight = 5.0

	// exaggerate the height of the cubes by multiplying them by this number.
	HeightMultiplier = 2

	// Oversize the blocks slightly to allow them to be merged together more cleanly.
	OverlapFactor = 1.2 //e.g. 1.2 = 120%
)

var (
	inputType = flag.String("input-type", "decimal", "Identifies what the input is. Can be decimal or text")
	shape     = flag.String("shapes-", "square", "What shape to use for the extrusions (circle or square)")
)

func init() {
	flag.Parse()
}

func main() {
	var input []int
	switch *inputType {
	case "decimal":
		input = mustGetNumbersFromDecimalNumber(os.Stdin)
	case "text":
		input = mustGetNumbersFromLetters(os.Stdin)
	default:
		panic(fmt.Sprintf("unknown input type: %s", *inputType))
	}
	render(input, *shape)
}

func render(heights []int, shape string) {
	grid := make([][]sdf.SDF3, 0)
	for i := 0; i < LineLength*NumLines; i++ {

		if len(grid) > NumLines {
			break
		}
		if len(grid) == 0 || len(grid[len(grid)-1]) == LineLength {
			grid = append(grid, []sdf.SDF3{})
		}

		var height = 0.0
		if len(heights) > i {
			height = float64(heights[i])
		}
		height = height * HeightMultiplier

		var s sdf.SDF2
		switch shape {
		case "square":
			// create the cube
			s = sdf.Box2D(
				sdf.V2{
					X: CubeSize * OverlapFactor,
					Y: CubeSize * OverlapFactor,
				},
				0,
			)
		case "circle":
			// create the cube
			s = sdf.Circle2D((CubeSize / 2) * OverlapFactor)
		default:
			panic(fmt.Sprintf("unknown shape: %s", shape))
		}
		c := sdf.Extrude3D(
			s,
			// make the cubes extend half way into the base to make the union cleaner.
			height+(BaseHeight/2),
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

	sdf.RenderSTL(base, 200, "result.stl")
}

func mustGetNumbersFromDecimalNumber(f io.Reader) []int {

	// todo: Should just read as much as needed rather than everything.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	numbers := []int{}

	beforeDecimalPoint := true
	for _, char := range string(b) {
		if string(char) == "." {
			beforeDecimalPoint = false
			continue
		}
		intVal, err := strconv.Atoi(string(char))
		if err != nil {
			panic(err)
		}
		if beforeDecimalPoint {
			numbers = append(numbers, intVal*10)
		} else {
			numbers = append(numbers, intVal)
		}
	}
	return numbers
}

func mustGetNumbersFromLetters(f io.Reader) []int {

	// todo: Should just read as much as needed rather than everything.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	numbers := make([]int, 0)
	charLookup := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16, "q": 17, "r": 18, "s": 19, "t": 20, "u": 21, "v": 22, "w": 23, "x": 24, "y": 25, "z": 26}
	for _, c := range string(b) {
		char := strings.ToLower(string(c))
		charNum, ok := charLookup[char]
		if !ok {
			if char == " " {
				numbers = append(numbers, 0)
			}
			continue
		}
		numbers = append(numbers, charNum)
	}
	return numbers
}
