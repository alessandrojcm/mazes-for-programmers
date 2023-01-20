package algorithms

import (
	"github.com/gookit/goutil/arrutil"
	"log"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

func findCell(arr []*mfp.Cell, target *mfp.Cell) (index int) {
	index = -1
	for i, cell := range arr {
		if cell == target {
			index = i
			return
		}
	}
	return
}

func filterCell(a []*mfp.Cell, target *mfp.Cell) []*mfp.Cell {
	n := 0
	for _, x := range a {
		if x != target {
			a[n] = x
			n++
		}
	}
	return a[:n]
}

func Wilson(grid grids.BaseGridHandler, cutOffPoint int) {
	log.Printf("starting wilson run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "wilson run")
	var unvisited []*mfp.Cell
	count := 0
	for cell := range grid.EachCell() {
		unvisited = append(unvisited, cell)
	}
	first := arrutil.GetRandomOne[*mfp.Cell](unvisited)
	// removing from array
	filterCell(unvisited, first)
	for len(unvisited) > 0 {
		if cutOffPoint != -1 && count >= cutOffPoint {
			return
		}
		cell := arrutil.GetRandomOne(unvisited)
		path := []*mfp.Cell{cell}

		for findCell(unvisited, cell) != -1 {
			cell = arrutil.GetRandomOne(cell.Neighbors())
			position := findCell(path, cell)
			if position != -1 {
				path = path[0:position]
			} else {
				path = append(path, cell)
			}

			for i := 0; i < len(path)-1; i++ {
				count++
				path[i].Link(path[i+1], true)
				unvisited = filterCell(unvisited, path[i])
			}
		}
	}
}
