package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp/grids"
	"os"
)

// showCmd -- sub command cli
var showCmd = &cobra.Command{
	Use:       "show",
	Short:     "Renders the maze to a window",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)

		if longestPath {
			grid, _ := builder.BuildGridWithDistanceRenderer(validColors[backgroundCol])
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
			// solve for start & end
			grid, _ := builder.BuildGridWithDistanceRenderer(validColors[backgroundCol])
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
		rl.ClearWindowState(rl.FlagWindowHidden)
		rl.SetWindowTitle(name)
		rl.SetTargetFPS(60)
		for !rl.WindowShouldClose() {
			rl.BeginDrawing()
			rl.DrawTexture(target.Texture, 0, 0, rl.White)
			rl.EndDrawing()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		postRenderCleanup()
	},
}

// animateCmd -- subcommand cmd for animating maze-solving
var animateCmd = &cobra.Command{
	Use:       "animate",
	Short:     "Animates the path rendering",
	ValidArgs: validArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var _ string
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		builder := grids.NewBuilder(rows, columns)
		grid, _ := builder.BuildAnimatableGrid(validColors[backgroundCol])

		if longestPath {
			n, solution, err := handleLongestPath(grid, handleAlgorithms(cmd, args, grid))
			_ = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// solve for start & end
			n, solution, err := handlePathSolve(grid, handleAlgorithms(cmd, args, grid))
			_ = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		} else {
			// if not anything set, solve for max path
			// solve for start & end
			startCell, endCell = "0x0", fmt.Sprintf("%dx%d", rows-1, columns-2)
			n, solution, err := handlePathSolve(grid, handleAlgorithms(cmd, args, grid))
			_ = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		}
		grid.ShowAnimation(cellSizes, thickness)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		postRenderCleanup()
	},
}

func InitShow(cmd *cobra.Command) {
	addRenderingFlags(showCmd)
	addSolvingFlags(showCmd)

	addSolvingFlags(animateCmd)
	addRenderingFlags(animateCmd)

	cmd.AddCommand(showCmd)
	cmd.AddCommand(animateCmd)
}
