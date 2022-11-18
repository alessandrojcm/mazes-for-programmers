package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mazes-for-programmers/commands"
	mfp2 "mazes-for-programmers/mfp/algorithms"
	"os"
	"runtime"
)

var debug bool
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
			rl.SetTraceLog(rl.LogNone)
			log.SetOutput(io.Discard)
		}
	},
	Aliases: []string{"mfp"},
}

// TODO: print weights of the cells for the show command(s)
func main() {
	// Raylib uses OpenGL and OpenGL expects every
	// call to be main on a single thread
	// so block the thread to avoid crashes
	runtime.LockOSThread()

	rootCmd.PersistentFlags().IntP("rows", "r", 4, "number or rows for the maze")
	rootCmd.PersistentFlags().IntP("columns", "c", 4, "number of columns for the maze")
	rootCmd.PersistentFlags().StringP("bias", "b", "", "set the bias for the algorithm, options are: "+mfp2.SouthAndWest+", "+
		", "+mfp2.NorthAndWest+
		mfp2.SouthAndEast)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "show debug grid")
	commands.InitPrint(rootCmd)
	commands.InitShow(rootCmd)
	commands.InitExport(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
