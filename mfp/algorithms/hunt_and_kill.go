package algorithms

import (
	"github.com/gookit/goutil/arrutil"
	"log"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

func HuntAndKill(grid grids.BaseGridHandler, cutOffPoint int) {
	log.Printf("starting hunt & kill tree run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "hunt & kill")
	count := 0
	current, _ := grid.RandomCell()
	for current != nil {
		if cutOffPoint != -1 && count >= cutOffPoint {
			return
		}
		unvisitedNeighbors := arrutil.TakeWhile(current.Neighbors(), func(a *mfp.Cell) bool {
			return len(a.Links()) == 0
		})

		if len(unvisitedNeighbors) > 0 {
			neighbor := arrutil.GetRandomOne[*mfp.Cell](unvisitedNeighbors)
			current.Link(neighbor, true)
			current = neighbor
		} else {
			current = nil
			for cell := range grid.EachCell() {
				visitedNeighbors := arrutil.TakeWhile(cell.Neighbors(), func(a *mfp.Cell) bool {
					for _, link := range a.Links() {
						if link != nil {
							return true
						}
					}
					return false
				})
				if len(cell.Links()) == 0 && len(visitedNeighbors) > 0 {
					count++
					current = cell
					neighbor := arrutil.GetRandomOne[*mfp.Cell](visitedNeighbors)
					current.Link(neighbor, true)
				}
			}
		}
	}
}
