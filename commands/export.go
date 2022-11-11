package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
)

var exportCmd = &cobra.Command{
	Use:       "export",
	Short:     "Exports the maze to an PNG image",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)
		grid, _ := builder.BuildGridLineRenderer()
		name := handleAlgorithms(cmd, args, grid)

		target := grid.ToTexture(cellSizes, thickness)

		defer rl.UnloadTexture(target.Texture)
		img := rl.LoadImageFromTexture(target.Texture)
		rl.ImageFlipVertical(*&img)
		defer rl.UnloadImage(img)
		rl.ExportImage(*img, fmt.Sprintf("%s.png", name))
		defer rl.UnloadRenderTexture(*target)
		defer rl.CloseWindow()
	},
}

func InitExport(cmd *cobra.Command) {
	addFlags(exportCmd)
	cmd.AddCommand(exportCmd)
}
