package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"mfp/commands"
	mfp2 "mfp/mfp"
	"os"
	"runtime"
)

// TODO: refactor code repetition on commands
// TODO: add flag to print ascii (make sure to require at least one of the showing flags)
// TODO: add flags to specify the start & end cell (maybe --solve-from and --solve-to or --solve=from-to idk)
// TODO: add mode to paint background with weight color
// TODO: refactor file structure
// TODO: add flag to paint the longest path (override --distance and require either longest or from-to solved)
// TODO: add loggin
// TODO: add timer
func main() {
	// Raylib uses OpenGL and OpenGL expects every
	// call to be main on a single thread
	// so block the thread to avoid crashes
	runtime.LockOSThread()
	var debug bool
	var rootCmd = &cobra.Command{Use: "render2d", Aliases: []string{"r2d"}, Run: func(cmd *cobra.Command, args []string) {
		if debug {
			os.Setenv("debug", "True")
		}
	}}
	rootCmd.PersistentFlags().IntP("rows", "r", 4, "number or rows for the maze")
	rootCmd.PersistentFlags().IntP("columns", "c", 4, "number of columns for the maze")
	rootCmd.PersistentFlags().StringP("bias", "b", "", "set the bias for the algorithm, options are: "+mfp2.SouthAndWest+", "+
		", "+mfp2.NorthAndWest+
		mfp2.SouthAndEast)
	rootCmd.PersistentFlags().BoolP("export", "e", false, "export the maze to an image")
	rootCmd.PersistentFlags().IntP("cellsize", "s", 10, "sets the size of the cells")
	rootCmd.PersistentFlags().IntP("thickness", "w", 1, "sets the thickness of the walls for the exported images")
	rootCmd.PersistentFlags().BoolP("tiles", "t", false, "colour the tiles")
	rootCmd.PersistentFlags().BoolP("show", "o", false, "show result in window")
	rootCmd.PersistentFlags().BoolP("distance", "i", false, "render the distance value of the cells")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "show debug grid")
	rootCmd.AddCommand(commands.Sidewinder)
	rootCmd.AddCommand(commands.BinaryTree)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
