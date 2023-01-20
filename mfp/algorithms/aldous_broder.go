package algorithms

import (
	"github.com/gookit/goutil/arrutil"
	"log"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

// AldousBroder -- an implementation of the Aldous-Broder algorithm
func AldousBroder(grid grids.BaseGridHandler, cutOffPoint int) {
	log.Printf("starting aldous-broder tree run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "aldous broder")
	cell, err := grid.RandomCell()
	unvisited, count := grid.Size()-1, 0
	if err != nil {
		panic(err)
	}

	for unvisited > 0 {
		if cutOffPoint != -1 && count >= cutOffPoint {
			return
		}
		neighbor := arrutil.GetRandomOne[*mfp.Cell](cell.Neighbors())

		if len(neighbor.Links()) == 0 {
			err = cell.Link(neighbor, true)
			if err != nil {
				panic(err)
			}
			unvisited -= 1
			log.Printf("%d cells to go", unvisited)
			count++
		}
		cell = neighbor
	}
}
