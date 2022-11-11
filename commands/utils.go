package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/algorithms"
	"mazes-for-programmers/mfp/grids"
	"os"
	"strconv"
	"strings"
)

var validArgs = []string{"sidewinder", "binarytree"}
var target *rl.RenderTexture2D

func handleAlgorithms(cmd *cobra.Command, args []string, grid grids.BaseGridHandler) (name string) {
	rows, _ := cmd.Flags().GetInt("rows")
	columns, _ := cmd.Flags().GetInt("columns")
	bias, _ := cmd.Flags().GetString("bias")

	if args[0] != "binarytree" && len(bias) > 0 {
		fmt.Fprintln(os.Stderr, "bias is only implemented for binarytree")
		os.Exit(1)
	}

	switch args[0] {
	case "binarytree":
		name = fmt.Sprintf("%s %vx%v maze with bias %s", "binarytree", rows, columns, bias)
		algorithms.BinaryTree(grid, bias)
	case "sidewinder":
		name = fmt.Sprintf("%s %vx%v maze with bias %s", "sidewinder", rows, columns, bias)
		algorithms.SideWinder(grid)
	}
	return
}

func handlePathSolve(grid grids.BaseGridHandler, name string) (string, mfp.Distance, error) {
	normalizedStart, normalizedEnd := strings.Split(strings.ToLower(startCell), "x"), strings.Split(strings.ToLower(endCell), "x")
	// get start & end cell
	startRow, _ := strconv.Atoi(normalizedStart[0])
	startCol, _ := strconv.Atoi(normalizedStart[1])
	endRow, _ := strconv.Atoi(normalizedEnd[0])
	endCol, _ := strconv.Atoi(normalizedEnd[1])
	// Calculate path
	start, err := grid.CellAt(startRow, startCol)
	if err != nil {
		return name, mfp.Distance{}, nil
	}
	end, err := grid.CellAt(endRow, endCol)
	if err != nil {
		return name, mfp.Distance{}, nil
	}
	distances := start.Distances()
	return fmt.Sprintf("%s from %s to %s", name, startCell, endCell), distances.PathTo(end), nil
}

func handleLongestPath(grid grids.BaseGridHandler, name string) (string, mfp.Distance, error) {
	start, err := grid.CellAt(0, 0)
	d := start.Distances()
	newStart, _ := d.Max()
	newDistances := newStart.Distances()
	goal, _ := newDistances.Max()
	if err != nil {
		return name, newDistances.PathTo(goal), err
	}
	return fmt.Sprintf("%slongest path", name), newDistances.PathTo(goal), nil
}

func postRenderCleanup() {
	rl.UnloadTexture(target.Texture)
	rl.UnloadRenderTexture(*target)
	rl.CloseWindow()
}
