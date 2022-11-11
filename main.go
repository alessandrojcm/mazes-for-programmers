package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"mazes-for-programmers/commands"
	mfp2 "mazes-for-programmers/mfp/algorithms"
	"os"
	"runtime"
)

// TODO: add mode to paint background with weight color
// TODO: add logging
// TODO: add timer
// TODO: add flags to specify the start & end cell (maybe --solve-from and --solve-to or --solve=from-to idk)
// TODO: add flag to paint the longest path (override --distance and require either longest or from-to solved)
func main() {
	// Raylib uses OpenGL and OpenGL expects every
	// call to be main on a single thread
	// so block the thread to avoid crashes
	runtime.LockOSThread()
	var debug bool
	var rootCmd = &cobra.Command{Use: "mfp", Aliases: []string{"mfp"}, Run: func(cmd *cobra.Command, args []string) {
		if debug {
			os.Setenv("DEBUG", "True")
		}
	}}
	rootCmd.PersistentFlags().IntP("rows", "r", 4, "number or rows for the maze")
	rootCmd.PersistentFlags().IntP("columns", "c", 4, "number of columns for the maze")
	rootCmd.PersistentFlags().StringP("bias", "b", "", "set the bias for the algorithm, options are: "+mfp2.SouthAndWest+", "+
		", "+mfp2.NorthAndWest+
		mfp2.SouthAndEast)
	//rootCmd.PersistentFlags().BoolP("distance", "i", false, "render the distance value of the cells")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "show debug grid")
	commands.InitPrint(rootCmd)
	commands.InitShow(rootCmd)
	commands.InitExport(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
