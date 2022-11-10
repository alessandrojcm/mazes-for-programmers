package grids

import (
	"mazes-for-programmers/mfp"
	"strconv"
)

type DistanceGrid struct {
	*ASCIIGrid
	Distances mfp.Distance
}

func (grid DistanceGrid) String() string {
	return gridToString(grid)
}

func (grid DistanceGrid) ContentsOf(cell *mfp.Cell) string {
	distance, isOk := grid.Distances.Cells[cell]
	if !isOk {
		return " "
	}
	return strconv.FormatInt(int64(distance), 36)
}
