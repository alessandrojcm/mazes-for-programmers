package algorithms

import (
	"log"
	"math/rand"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

// SideWinder -- A sidewinder algorithm implementation
func SideWinder(grid grids.BaseGridHandler) {
	log.Printf("starting sidewinder run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "sidewinder run")
	var err error
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for row := range grid.EachRow() {
		run := make([]*mfp.Cell, 0, len(row))

		for _, cell := range row {
			if cell == nil {
				continue
			}
			run = append(run, cell)
			atEasternBoundary, atNorthernBoundary := cell.East == nil, cell.North == nil
			shouldCloseOut := atEasternBoundary || (!atNorthernBoundary && r.Intn(2) == 0)

			if shouldCloseOut {
				sample := r.Intn(len(run))
				member := run[sample]
				if member.North != nil {
					err = member.Link(member.North, true)
				}
				run = make([]*mfp.Cell, 0, len(row))
			} else if cell.East != nil {
				err = cell.Link(cell.East, true)
			}
			if err != nil {
				panic(err)
			}
		}
	}
}
