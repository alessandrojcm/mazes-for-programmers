package mfp

import (
	"strconv"
)

type DistanceGrid struct {
	*ASCIIGrid
}

func (grid *DistanceGrid) ContentsOf(cell *Cell) string {
	distance, isOk := cell.Distances().cells[&cell]
	if !isOk {
		return " "
	}
	return strconv.FormatInt(int64(distance), 36)
}
