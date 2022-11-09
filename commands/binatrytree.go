package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mfp/flags"
	"mfp/mfp"
	"time"
)

var BinaryTree = &cobra.Command{
	Use:     "binarytree",
	Aliases: []string{"bt"},
	Short:   "Renders a maze using the sidewinder algorithm",
	Long:    "The sidewinder algorithm tries to group adjacent cells together before carving out a passage north of them.",
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		cellSize, _ := cmd.Flags().GetInt("cellsize")
		wallThickness, _ := cmd.Flags().GetInt("thickness")
		export, _ := cmd.Flags().GetBool("export")
		colorTiles, _ := cmd.Flags().GetBool("tiles")
		show, _ := cmd.Flags().GetBool("show")
		bias, _ := cmd.Flags().GetString("bias")
		distance, _ := cmd.Flags().GetBool("distance")
		builder := mfp.NewBuilder(rows, columns)
		if distance {
			grid, err := builder.BuildGridWithDistance()
			mfp.BinaryTree(grid, bias)
			start, err := grid.CellAt(0, 0)

			distances := start.Distances()
			grid.Distances = distances

			fmt.Printf("Printing distance %s %vx%v maze with %s bias\n", "binary tree", rows, columns, bias)
			fmt.Println(grid)
			goal, err := grid.CellAt(rows-1, columns-1)
			grid.Distances = distances.PathTo(goal)
			fmt.Printf("Printing solved %s %vx%v maze with %s bias\n", "binary tree", rows, columns, bias)
			fmt.Println(grid)

			if err != nil {
				panic(err)
			}
		} else {
			grid, err := builder.BuildASCIIGrid()
			if err != nil {
				panic(err)
			}
			mfp.BinaryTree(grid, bias)
			fmt.Printf("Printing %s %vx%v maze with %s bias\n", "sidewinder", rows, columns, bias)
			fmt.Println(grid)
		}

		if export || show {
			var target *rl.RenderTexture2D
			var err error

			var name string
			if distance {
				name = fmt.Sprintf("%s-%vrowX%vcol-%s-%v-with-distance.png", "sidewinder", rows, columns, bias, time.Now().UnixNano())

			} else {
				name = fmt.Sprintf("%s-%vrowX%vcol-%s-%v.png", "sidewinder", rows, columns, bias, time.Now().UnixNano())
			}
			if colorTiles {
				rendererGrid, _ := builder.BuildGridTiledRenderer()
				target = rendererGrid.ToTexture(cellSize, wallThickness)
			} else {
				rendererGrid, _ := builder.BuildGridLineRenderer()
				target = rendererGrid.ToTexture(cellSize, wallThickness)
			}
			if err != nil {
				panic(err)
			}
			if export {
				flags.Export(target, name)
			}
			if show {
				flags.Show(target, name)
			}
			defer rl.UnloadRenderTexture(*target)
			defer rl.CloseWindow()
		}
	},
}
