package grids

import (
	"mazes-for-programmers/mfp"
	"strconv"
)

// DistanceGrid -- A grid that can measure distances between cells
type DistanceGrid struct {
	*ASCIIGrid
	Distances mfp.Distance
}

func (grid DistanceGrid) String() string {
	return gridToString(grid)
}

// ContentsOf -- Prints the weight of the cell in base36
func (grid DistanceGrid) ContentsOf(cell *mfp.Cell) string {
	distance, isOk := grid.Distances.Cells[cell]
	if !isOk {
		return " "
	}
	return strconv.FormatInt(int64(distance), 36)
}
