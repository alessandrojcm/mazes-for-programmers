package algorithms

import (
	"github.com/gookit/goutil/arrutil"
	"log"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

func HuntAndKill(grid grids.BaseGridHandler) {
	log.Printf("starting hunt & kill tree run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "hunt & kill")
	current, _ := grid.RandomCell()
	for current != nil {
		unvisitedNeighbors := arrutil.TakeWhile(current.Neighbors(), func(a any) bool {
			return len(a.(*mfp.Cell).Links()) == 0
		}).([]*mfp.Cell)

		if len(unvisitedNeighbors) > 0 {
			neighbor := arrutil.GetRandomOne[*mfp.Cell](unvisitedNeighbors)
			current.Link(neighbor, true)
			current = neighbor
		} else {
			current = nil
			for cell := range grid.EachCell() {
				visitedNeighbors := arrutil.TakeWhile(cell.Neighbors(), func(a any) bool {
					for _, link := range a.(*mfp.Cell).Links() {
						if link != nil {
							return true
						}
					}
					return false
				}).([]*mfp.Cell)
				if len(cell.Links()) == 0 && len(visitedNeighbors) > 0 {
					current = cell
					neighbor := arrutil.GetRandomOne[*mfp.Cell](visitedNeighbors)
					current.Link(neighbor, true)
				}
			}
		}
	}
}
