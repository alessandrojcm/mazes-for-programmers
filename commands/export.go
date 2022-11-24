package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
	"os"
)

// exportCmd sub command line utility to handle image exports
var exportCmd = &cobra.Command{
	Use:       "export",
	Short:     "Exports the maze to an PNG image",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)
		var name string

		if longestPath {
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			n, solution, err := handleLongestPath(grid, handleAlgorithms(cmd, args, grid))
			name = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness)
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// grid with path solving
			//  for start & end
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			n, solution, err := handlePathSolve(grid, handleAlgorithms(cmd, args, grid))
			name = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness)
		} else {
			// Normal grid
			grid, _ := builder.BuildGridLineRenderer()
			name = handleAlgorithms(cmd, args, grid)
			target = grid.ToTexture(cellSizes, thickness)
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
