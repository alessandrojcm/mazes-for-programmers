package mfp

import (
	"math/rand"
	"time"
)

func SideWinder(grid BaseGridHandler) {
	var err error
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for row := range grid.EachRow() {
		run := make([]*Cell, 0, len(row))

		for _, cell := range row {
			if cell == nil {
				continue
			}
			run = append(run, cell)
			atEasternBoundary, atNorthernBoundary := cell.east == nil, cell.north == nil
			shouldCloseOut := atEasternBoundary || (!atNorthernBoundary && r.Intn(2) == 0)

			if shouldCloseOut {
				sample := r.Intn(len(run))
				member := run[sample]
				if member.north != nil {
					err = member.Link(member.north, true)
				}
				run = make([]*Cell, 0, len(row))
			} else if cell.east != nil {
				err = cell.Link(cell.east, true)
			}
			if err != nil {
				panic(err)
			}
		}
	}
}
