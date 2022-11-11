package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/algorithms"
	"mazes-for-programmers/mfp/grids"
	"os"
)

var validArgs = []string{"sidewinder", "binarytree"}

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
		name = fmt.Sprintf("Printing %s %vx%v maze with bias %s\n", "binarytree", rows, columns, bias)
		algorithms.BinaryTree(grid, bias)
	case "sidewinder":
		name = fmt.Sprintf("Printing %s %vx%v maze with bias %s\n", "sidewinder", rows, columns, bias)
		algorithms.SideWinder(grid)
	}
	return
}
