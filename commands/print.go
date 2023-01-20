package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"os"
)

// printCmd -- sub cmd utility to print an ASCII maze
var printCmd = &cobra.Command{
	Use:       "print",
	Short:     "Prints a maze wih ASCII characters",
	Aliases:   []string{"p"},
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.RangeArgs(1, len(validArgs)-1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)
		var name string
		var solution mfp.Distance
		var err error

		// print longest path
		if longestPath {
			grid, _ := builder.BuildGridWithDistance()
			name, err = handleAlgorithms(cmd, args, grid)
			name, solution, err = handleLongestPath(grid, name)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			fmt.Println(name, "\n", grid)
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// solve for start & end
			grid, _ := builder.BuildGridWithDistance()
			name, err = handleAlgorithms(cmd, args, grid)
			name, solution, err = handlePathSolve(grid, name)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			fmt.Println(name, "\n", grid)
		} else if spreadMiddle {
			grid, _ := builder.BuildGridWithDistance()
			middle, err := grid.CellAt(rows/2, columns/2)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = middle.Distances()
		} else {
			// print normal maze
			grid, _ := builder.BuildASCIIGrid()
			name, err = handleAlgorithms(cmd, args, grid)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			fmt.Println(name, "\n", grid)
		}
	},
}

// distancesCmd -- sub command to print the weight of each cell
var distancesCmd = &cobra.Command{
	Use:     "distances",
	Short:   "Prints the distance value of every cell",
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)

		grid, _ := builder.BuildGridWithDistance()
		name, err := handleAlgorithms(cmd, args, grid)
		start, _ := grid.CellAt(0, 0)
		distances := start.Distances()
		grid.Distances = distances
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(-1)
		}

		fmt.Println(name, "\n", grid)
	},
}

func InitPrint(cmd *cobra.Command) {
	// We do not add the flags to distance since that command renders all the values for the cells
	addSolvingFlags(printCmd)
	printCmd.AddCommand(distancesCmd)
	cmd.AddCommand(printCmd)
}
