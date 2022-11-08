package commands

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"mfp/flags"
	"mfp/mfp"
	"time"
)

var Sidewinder = &cobra.Command{
	Use:     "sidewinder",
	Aliases: []string{"sd"},
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
		builder := mfp.NewBuilder(rows, columns)
		grid, err := builder.BuildASCIIGrid()
		if err != nil {
			panic(err)
		}
		mfp.SideWinder(grid)
		fmt.Printf("Printing %s %vx%v maze with %s bias\n", "sidewinder", rows, columns, bias)
		fmt.Println(grid)
		if export || show {
			rendererGrid, err := builder.BuildRenderGrid()
			if err != nil {
				panic(err)
			}
			name := fmt.Sprintf("%s-%vrowX%vcol-%s-%v.png", "sidewinder", rows, columns, bias, time.Now().UnixNano())
			target := rendererGrid.ToTexture(cellSize, wallThickness, colorTiles)

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
