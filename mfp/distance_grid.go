package mfp

import (
	"strconv"
)

type DistanceGrid struct {
	*ASCIIGrid
	Distances Distance
}

func (grid DistanceGrid) String() string {
	return gridToString(grid)
}

func (grid DistanceGrid) ContentsOf(cell *Cell) string {
	distance, isOk := grid.Distances.cells[cell]
	if !isOk {
		return " "
	}
	return strconv.FormatInt(int64(distance), 36)
}
