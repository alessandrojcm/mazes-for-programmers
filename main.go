package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mazes-for-programmers/commands"
	mfp "mazes-for-programmers/mfp/algorithms"
	"os"
	"strings"
)

var debug bool

// root command line
var rootCmd = &cobra.Command{
	Use: "mfp",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		columns, _ := cmd.Flags().GetInt("columns")
		if rows <= 0 || columns <= 0 {
			cmd.PrintErrln("rows and columns need to be greater than 0.")
			os.Exit(-1)
		}
		if debug {
			os.Setenv("DEBUG", "True")
			log.SetFlags(log.Ldate)
		} else {
			rl.SetTraceLogLevel(rl.LogNone)
			log.SetOutput(io.Discard)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
        cmd.Println("Mazes for Programmers (mfp) is a Go implementation of the Mazes For Programmers book by Jamis buck.\n")
		cmd.Println("The maze supported maze algorithms are:\n1) Aldous Broder\n2) Binary Tree\n3) Hunt & kill\n4) Sidewinder\n5) Wilson\nYou can also use mixup followed by two or more algorithms to combine them.")
    },

	Aliases: []string{"mfp"},
}

// TODO: fix weight fot the ASCII version (it overflows the cells)
// TODO: spread middle seems fishy, check it out
// TODO: prettify text?
func main() {
	rootCmd.PersistentFlags().IntP("rows", "r", 4, "number or rows for the maze")
	rootCmd.PersistentFlags().IntP("columns", "c", 4, "number of columns for the maze")
	rootCmd.PersistentFlags().StringP("bias", "b", "", fmt.Sprintf("set the bias for the algorithm, options are: %s", strings.Join([]string{mfp.SouthAndWest,
		mfp.NorthAndWest,
		mfp.SouthAndEast}, ", ")))
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "show debug grid")
	commands.InitShow(rootCmd)
	commands.InitExport(rootCmd)
	commands.InitPrint(rootCmd)
	commands.InitDeadEnds(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
