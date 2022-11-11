package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
)

var cellSizes, thickness int
var startCell, endCell string

func addRenderingFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&cellSizes, "cellsize", "s", 10, "sets the size of the cells")
	cmd.Flags().IntVarP(&thickness, "thickness", "w", 1, "sets the thickness of the walls for the exported images")
}

func addSolvingFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().StringVarP(&startCell, "solve-from", "f", "", "Set the starting cell to solve the maze from in the RowxColumn format")
		cmd.Flags().StringVarP(&endCell, "solve-to", "t", "", "Set the ending cell to solve the maze to in the RowxColumn format")
		cmd.MarkFlagsRequiredTogether("solve-from", "solve-to")

		cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			rows, _ := cmd.Flags().GetInt("rows")
			columns, _ := cmd.Flags().GetInt("columns")
			validRange := regexp.MustCompile(fmt.Sprintf("(?mi)[0-%d]{1,%d}x[0-%d]{1,%d}", rows-1, len(strconv.Itoa(rows-1)), columns-1, len(strconv.Itoa(columns-1))))
			if len(startCell) == 0 || len(endCell) == 0 {
				return
			}
			if !validRange.MatchString(startCell) {
				cmd.PrintErr("Starting cell out of range\n")
				os.Exit(-1)
			}
			if !validRange.MatchString(endCell) {
				cmd.PrintErr("End cell out of range\n")
				os.Exit(-1)
			}
			if startCell == endCell {
				cmd.PrintErr("Start and end cells cannot be the same\n")
				os.Exit(-1)
			}
		}
	}
}
