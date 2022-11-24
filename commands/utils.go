package commands

import (
	"errors"
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

// handleAlgorithms -- receives a grid and applies the corresponding algorithm to it, also returns the file/window name
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

// parseCellExpression -- validates expressions for the rowxcol input an returns the row and the col
func parseCellExpression(cellExpr string) ([]int, error) {
	err := errors.New("malformed cell expression, make sure is in the form RowXColumn")
	if !strings.ContainsAny(cellExpr, "xX") {
		return []int{}, err
	}
	cell := strings.Split(strings.ToLower(cellExpr), "x")
	if len(cell) != 2 {
		return []int{}, err
	}
	var row, col int
	row, err = strconv.Atoi(cell[0])
	if err != nil {
		return []int{}, err
	}
	col, err = strconv.Atoi(cell[1])
	return []int{row, col}, err
}

// handlePathSolve -- parses the input expression from the from-to path
func handlePathSolve(grid grids.BaseGridHandler, name string) (string, mfp.Distance, error) {
	normalizedStart, _ := parseCellExpression(startCell)
	normalizedEnd, _ := parseCellExpression(endCell)
	// Calculate path
	start, err := grid.CellAt(normalizedStart[0], normalizedStart[1])
	if err != nil {
		return name, mfp.Distance{}, nil
	}
	end, err := grid.CellAt(normalizedEnd[0], normalizedEnd[1])
	if err != nil {
		return name, mfp.Distance{}, nil
	}
	distances := start.Distances()
	return fmt.Sprintf("%s from %s to %s", name, startCell, endCell), distances.PathTo(end), nil
}

// handleLongestPath -- parses computing the largest path
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

// postRenderCleanup -- finalizes the opengl context
func postRenderCleanup() {
	if target != nil {
		rl.UnloadTexture(target.Texture)
		rl.UnloadRenderTexture(*target)
	}
	rl.CloseWindow()
}
