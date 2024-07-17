package commands

import (
	"errors"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gookit/goutil/arrutil"
	"github.com/spf13/cobra"
	"math"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/algorithms"
	"mazes-for-programmers/mfp/grids"
	"os"
	"sort"
	"strconv"
	"strings"
)

var validArgs = []string{"sidewinder", "binarytree", "aldous-broder", "wilson", "hunt-and-kill", "mixup"}
var target *rl.RenderTexture2D

// handleAlgorithms -- receives a grid and applies the corresponding algorithm to it, also returns the file/window name
func handleAlgorithms(cmd *cobra.Command, args []string, grid grids.BaseGridHandler) (name string, err error) {
	rows, _ := cmd.Flags().GetInt("rows")
	columns, _ := cmd.Flags().GetInt("columns")
	bias, _ := cmd.Flags().GetString("bias")
	cutOffPoint := -1

	if args[0] != "binarytree" && len(bias) > 0 {
		fmt.Fprintln(os.Stderr, "bias is only implemented for binarytree")
		os.Exit(1)
	}

	runAlgorithm := func(algo string) {
		switch algo {
		case "binarytree":
			name = fmt.Sprintf("%s %vx%v maze with bias %s", "binarytree", rows, columns, bias)
			algorithms.BinaryTree(grid, bias, cutOffPoint)
		case "sidewinder":
			name = fmt.Sprintf("%s %vx%v maze with bias %s", "sidewinder", rows, columns, bias)
			algorithms.SideWinder(grid, cutOffPoint)
		case "aldous-broder":
			name = fmt.Sprintf("#{aldous-broder} #{%v}x#[%v} maze", rows, columns)
			algorithms.AldousBroder(grid, cutOffPoint)
		case "wilson":
			name = fmt.Sprintf("#{wilson} #{%v}x#[%v} maze", rows, columns)
			algorithms.Wilson(grid, cutOffPoint)
		case "hunt-and-kill":
			name = fmt.Sprintf("#{hunt-and-kill} #{%v}x#[%v} maze", rows, columns)
			algorithms.HuntAndKill(grid, cutOffPoint)
		}
	}

	if args[0] != "mixup" {
		runAlgorithm(args[0])
		return
	}

	algs := args[1:]
	if len(algs) < 2 {
		return name, errors.New("no algorithms provided to mix up")
	}
	// for some reason, Aldous-Broder gets stuck if not used
	// as a first algorithm, thus we will order the array
	// to guarantee it is used first
	sort.Slice(algs, func(i, j int) bool {
		if algs[i] == "aldous-broder" {
			return true
		}
		return false
	})
	for _, alg := range algs {
		if !arrutil.Contains(validArgs, alg) {
			return name, errors.New("invalid algorithm passed to mixup")
		}
	}
	cutOffPoint = int(math.Ceil(float64(grid.Size() / len(algs))))
	for _, alg := range algs {
		runAlgorithm(alg)
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
