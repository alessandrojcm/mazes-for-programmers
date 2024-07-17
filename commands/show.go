package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"os"
)

// showCmd -- sub command cli to render the maze to an image
var showCmd = &cobra.Command{
	Use:       "show",
	Short:     "Renders the maze to a window",
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
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness, printWeights)
		} else if len(startCell) > 0 && len(endCell) > 0 {
			// solve for start & end
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			name, err = handleAlgorithms(cmd, args, grid)
			n, solution, err = handlePathSolve(grid, name)
			name = n
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			grid.Distances = solution
			// create texture
			target = grid.ToTexture(cellSizes, thickness, printWeights)
		} else if spreadMiddle {
			grid, _ := builder.BuildGridWithDistanceRenderer(rl.Color(backgroundCol))
			_, err = handleAlgorithms(cmd, args, grid)
			middle, err := grid.CellAt(rows/2, columns/2)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = middle.Distances()
			target = grid.ToTexture(cellSizes, thickness, printWeights)
		} else {
			// Normal grid
			grid, _ := builder.BuildGridLineRenderer()
			name, err = handleAlgorithms(cmd, args, grid)
			target = grid.ToTexture(cellSizes, thickness)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
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
	Args:      cobra.MatchAll(cobra.RangeArgs(1, len(validArgs)-1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var _, n string
		var solution mfp.Distance
		var err error
		var name string
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		breadcrumb, _ := cmd.Flags().GetBool("breadcrumb")
		builder := grids.NewBuilder(rows, columns)
		grid, _ := builder.BuildGridWithDistance()

		if longestPath {
			name, err = handleAlgorithms(cmd, args, grid)
			n, solution, err = handleLongestPath(grid, name)
			_ = n
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		} else if len(startCell) > 0 && len(endCell) > 0 {
			name, err = handleAlgorithms(cmd, args, grid)
			// solve for start & end
			n, solution, err = handlePathSolve(grid, name)
			_ = n
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		} else if spreadMiddle {
			_, err = handleAlgorithms(cmd, args, grid)
			middle, err := grid.CellAt(rows/2, columns/2)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = middle.Distances()
		} else {
			// if not anything set, solve for max path
			// solve for start & end
			startCell, endCell = "0x0", fmt.Sprintf("%dx%d", rows-1, columns-2)
			name, err = handleAlgorithms(cmd, args, grid)
			n, solution, err := handlePathSolve(grid, name)
			_ = n
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			grid.Distances = solution
		}
		if breadcrumb {
			renderer, _ := builder.BuildGridWithBreadCrumb(rl.Color(backgroundCol))
			renderer.ShowAnimation(cellSizes, thickness)
		} else {
			renderer, _ := builder.BuildAnimatableGrid(rl.Color(backgroundCol))
			renderer.ShowAnimation(cellSizes, thickness, printWeights)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		postRenderCleanup()
	},
}

func InitShow(cmd *cobra.Command) {
	addRenderingFlags(showCmd)
	addSolvingFlags(showCmd, animateCmd)

	addRenderingFlags(animateCmd)
	animateCmd.Flags().Bool("breadcrumb", false, "draw a breadcrumb line")

	cmd.AddCommand(showCmd)
	cmd.AddCommand(animateCmd)
	animateCmd.MarkFlagsOneRequired("print-weights", "solve-from", "solve-to", "longest-path")
	showCmd.MarkFlagsOneRequired("print-weights", "solve-from", "solve-to", "longest-path")
}
