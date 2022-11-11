package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
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

		grid, _ := builder.BuildASCIIGrid()
		name := handleAlgorithms(cmd, args, grid)
		fmt.Println(name, grid)
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

		fmt.Println(name, grid)
	},
}

func InitPrint(cmd *cobra.Command) {
	printCmd.AddCommand(distancesCmd)
	cmd.AddCommand(printCmd)
}
