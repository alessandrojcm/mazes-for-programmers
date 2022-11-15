package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
	"os"
)

var printCmd = &cobra.Command{
	Use:       "print",
	Short:     "Prints a maze wih ASCII characters",
	Aliases:   []string{"p"},
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)

		if longestPath {
			// solve for start & end
			grid, _ := builder.BuildGridWithDistance()
			name, solution, err := handleLongestPath(grid, handleAlgorithms(cmd, args, grid))
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			fmt.Println(name, "\n", grid)
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// solve for start & end
			grid, _ := builder.BuildGridWithDistance()
			name, solution, err := handlePathSolve(grid, handleAlgorithms(cmd, args, grid))
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			fmt.Println(name, "\n", grid)
		} else {
			grid, _ := builder.BuildASCIIGrid()
			name := handleAlgorithms(cmd, args, grid)
			fmt.Println(name, "\n", grid)
		}
	},
}

var distancesCmd = &cobra.Command{
	Use:     "distances",
	Short:   "Prints the distance value of every cell",
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)

		grid, _ := builder.BuildGridWithDistance()
		name := handleAlgorithms(cmd, args, grid)
		start, _ := grid.CellAt(0, 0)
		distances := start.Distances()
		grid.Distances = distances

		fmt.Println(name, "\n", grid)
	},
}

func InitPrint(cmd *cobra.Command) {
	// We do not add the flags to distance since that command renders all the values for the cells
	addSolvingFlags(printCmd)
	printCmd.AddCommand(distancesCmd)
	cmd.AddCommand(printCmd)
}
