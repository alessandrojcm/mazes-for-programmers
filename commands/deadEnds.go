package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "mazes-for-programmers/mfp/grids"
    "os"
    "github.com/jedib0t/go-pretty/v6/table"
)

// TODO: fix this
func calculateAlgorithmDeadEnd(tries int, builder *grids.GridBuilder, cmd *cobra.Command, args []string) int {
    var deadEndCount []int
    for i := 0; i < tries; i++ {
        grid, err := builder.BuildASCIIGrid()
        _, err = handleAlgorithms(cmd, args, grid)
        if err != nil {
            cmd.PrintErr(err)
            os.Exit(-1)
        }
        deadEndCount = append(deadEndCount, len(grid.DeadCells()))
    }
    totalDeadEnds := 0
    for _, deadEnd := range deadEndCount {
        totalDeadEnds += deadEnd
    }
    return totalDeadEnds / len(deadEndCount)
}

var deadEndValidArgs = []string{"sidewinder", "binarytree", "aldous-broder", "wilson", "hunt-and-kill", "all"}
var deadEndCmd = &cobra.Command{
    Use:   "dead-ends",
    Short: "Averages all the dead ends for each algorithms (or all of them)",
    Args:  cobra.MatchAll(cobra.RangeArgs(1, len(deadEndValidArgs)-1), cobra.OnlyValidArgs),
    Run: func(cmd *cobra.Command, args []string) {
        rows, _ := cmd.Flags().GetInt("rows")
        columns, _ := cmd.Flags().GetInt("columns")
        tries, _ := cmd.Flags().GetInt("tries")
        averages :=  make(map[string]int)
        builder := grids.NewBuilder(rows, columns)
        t := table.NewWriter()
        t.SetOutputMirror(cmd.OutOrStdout())
        t.AppendHeader(table.Row{"Algorithm", "Dead Ends", "Size", "%"})
        if args[0] == "all" {
            for _, algo := range deadEndValidArgs {
                if algo == "all" {
                    continue
                }
                cmd.Println("Running ", algo)

                averages[algo] = calculateAlgorithmDeadEnd(tries, builder, cmd, args)
            }
        } else {
			cmd.Println("Running for ", args[0])
			averages[args[0]] = calculateAlgorithmDeadEnd(tries, builder, cmd, args)
		}
		for k,v := range averages {
			percentage := (float32(v)*100.0) / float32(rows*columns)
            t.AppendRow(table.Row{k,v,rows*columns,percentage})
            t.AppendSeparator()
		}
		t.Render()
    },
}

func InitDeadEnds(cmd *cobra.Command) {
    deadEndCmd.Flags().Int("tries", 20, "Number of times the algorithm will be run")
    cmd.AddCommand(deadEndCmd)
}
