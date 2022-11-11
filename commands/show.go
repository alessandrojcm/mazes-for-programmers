package commands

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
)

var showCmd = &cobra.Command{
	Use:       "show",
	Short:     "Renders the maze to a window",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)
		grid, _ := builder.BuildGridLineRenderer()
		name := handleAlgorithms(cmd, args, grid)

		target := grid.ToTexture(cellSizes, thickness)

		rl.ClearWindowState(rl.FlagWindowHidden)
		rl.SetWindowTitle(name)
		defer rl.UnloadTexture(target.Texture)
		rl.SetTargetFPS(60)
		for !rl.WindowShouldClose() {
			rl.BeginDrawing()
			rl.DrawTexture(target.Texture, 0, 0, rl.White)
			rl.EndDrawing()
		}
		defer rl.UnloadRenderTexture(*target)
		defer rl.CloseWindow()
	},
}

func InitShow(cmd *cobra.Command) {
	addFlags(showCmd)

	cmd.AddCommand(showCmd)
}
