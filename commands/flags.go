package commands

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
)

type colorFlag rl.Color

var cellSizes, thickness int
var startCell, endCell string
var longestPath, spreadMiddle bool

var backgroundCol = colorFlag(rl.Green)

var validColors = map[string]rl.Color{
	"red":    rl.Red,
	"green":  rl.Green,
	"blue":   rl.Blue,
	"random": rl.Blank,
}

func (c *colorFlag) String() string {
	for k, v := range validColors {
		if v == rl.Color(*c) {
			return k
		}
	}
	return ""
}

func (c *colorFlag) Set(s string) error {
	parsedCol, isOk := validColors[s]
	if !isOk {
		return errors.New("invalid color")
	}
	*c = colorFlag(parsedCol)
	return nil
}

func (c *colorFlag) Type() string {
	return "color"
}

// addRendering-flags -- helper function to add all the flags for the "show" subcommand
func addRenderingFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&cellSizes, "cellsize", "s", 60, "sets the size of the cells")
	cmd.Flags().IntVarP(&thickness, "thickness", "w", 1, "sets the thickness of the walls for the exported images")
	cmd.Flags().VarP(&backgroundCol, "background-color", "l", "Set the background color to draw the distance grid with; options: red, green, blue and random.")
}

// addSolvingFlags -- helper function to add all the flags common to finding a solution for the maze
func addSolvingFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().StringVar(&startCell, "solve-from", "", "Set the starting cell to solve the maze from in the RowxColumn format")
		cmd.Flags().StringVar(&endCell, "solve-to", "", "Set the ending cell to solve the maze to in the RowxColumn format")
		cmd.Flags().BoolVarP(&longestPath, "longest-path", "p", false, "Draw the longest path in the maze")
		cmd.Flags().BoolVar(&spreadMiddle, "spread-middle", false, "Set the middle cell and start spreading from there")
		cmd.MarkFlagsRequiredTogether("solve-from", "solve-to")
		cmd.MarkFlagsMutuallyExclusive("solve-to", "longest-path")
		cmd.MarkFlagsMutuallyExclusive("solve-from", "longest-path")
		cmd.MarkFlagsMutuallyExclusive("spread-middle", "solve-from")
		cmd.MarkFlagsMutuallyExclusive("spread-middle", "solve-to")
		cmd.MarkFlagsMutuallyExclusive("spread-middle", "longest-path")

		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			rows, _ := cmd.Flags().GetInt("rows")
			columns, _ := cmd.Flags().GetInt("columns")
			if len(startCell) == 0 || len(endCell) == 0 {
				return nil
			}
			parsedStartCell, err := parseCellExpression(startCell)
			if err != nil {
				return errors.New("start cell: " + err.Error())
			}
			parsedEndCell, err := parseCellExpression(endCell)

			if startCell == endCell {
				return errors.New("Start and end cells cannot be the same\n")
			}
			if parsedStartCell[0] > rows-1 || parsedStartCell[0] < 0 || parsedStartCell[1] > columns-1 || parsedStartCell[1] < 0 {
				return errors.New("start cell out of range")
			}
			if parsedEndCell[0] > rows-1 || parsedEndCell[0] < 0 || parsedEndCell[1] > columns-1 || parsedEndCell[1] < 0 {
				return errors.New("start cell out of range")
			}
			return nil
		}
	}
}
