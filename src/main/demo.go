package main

import (
	"flag"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"mfp/src/mfp"
	"runtime"
	"time"
)

func main() {
	// Raylib uses OpenGL and OpenGL expects every
	// call to be main on a single thread
	// so block the thread to avoid crashes
	runtime.LockOSThread()
	var rows, columns, cellSize, wallThickness int
	var bias, algorithm string
	var export bool

	flag.IntVar(&rows, "rows", 4, "number or rows for the maze")
	flag.IntVar(&columns, "columns", 4, "number of columns for the maze")
	flag.StringVar(&bias, "bias", "", "set the bias for the algorithm, options are: "+mfp.SouthAndWest+", "+
		", "+mfp.NorthAndWest+
		mfp.SouthAndEast)
	flag.StringVar(&algorithm, "algorithm", "binarytree", "algorithm to use, options: binarytree, sidewinder")
	flag.BoolVar(&export, "export", false, "export the maze to an image")
	flag.IntVar(&cellSize, "cellsize", 10, "sets the size of the cells")
	flag.IntVar(&wallThickness, "thickness", 1, "sets the thickness of the walls for the exported images")
	flag.Parse()

	grid, err := mfp.NewGrid(rows, columns)
	if err != nil {
		panic(err)
	}
	switch algorithm {
	case "sidewinder":
		mfp.SideWinder(grid)
	default:
		mfp.BinaryTree(grid, bias)
	}
	fmt.Printf("Printing %s %vx%v maze with %s bias\n", algorithm, rows, columns, bias)
	fmt.Println(grid)
	if export {
		img := grid.ToImage(cellSize, wallThickness)
		defer rl.UnloadImage(img)
		rl.ExportImage(*img, fmt.Sprintf("%s-%vrowX%vcol-%s-%v.png", algorithm, rows, columns, bias, time.Now().UnixNano()))
	}
}
