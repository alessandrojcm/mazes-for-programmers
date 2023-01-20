package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"os"
)

// exportCmd sub command line utility to handle image exports
var exportCmd = &cobra.Command{
	Use:       "export",
	Short:     "Exports the maze to an PNG image",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.RangeArgs(1, len(validArgs)-1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var _, n string
		var solution mfp.Distance
		var err error
		var name string
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)

		if longestPath {
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			name, err = handleAlgorithms(cmd, args, grid)
			n, solution, err = handleLongestPath(grid, name)
			name = n
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness)
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// grid with path solving
			//  for start & end
			var err error
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			name, err = handleAlgorithms(cmd, args, grid)
			n, solution, err := handlePathSolve(grid, name)
			name = n
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness)
		} else {
			// Normal grid
			var err error
			grid, _ := builder.BuildGridLineRenderer()
			name, err = handleAlgorithms(cmd, args, grid)
			target = grid.ToTexture(cellSizes, thickness)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
		}
		img := rl.LoadImageFromTexture(target.Texture)
		rl.ImageFlipVertical(*&img)
		defer rl.UnloadImage(img)
		rl.ExportImage(*img, fmt.Sprintf("%s.png", name))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		postRenderCleanup()
	},
}

func InitExport(cmd *cobra.Command) {
	addRenderingFlags(exportCmd)
	addSolvingFlags(exportCmd)
	cmd.AddCommand(exportCmd)
}
