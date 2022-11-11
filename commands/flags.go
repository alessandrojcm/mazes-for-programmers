package commands

import "github.com/spf13/cobra"

var cellSizes, thickness int

func addFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&cellSizes, "cellsize", "s", 10, "sets the size of the cells")
	cmd.Flags().IntVarP(&thickness, "thickness", "w", 1, "sets the thickness of the walls for the exported images")
}
