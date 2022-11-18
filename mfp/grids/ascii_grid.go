package grids

import "mazes-for-programmers/mfp"

// ASCIIGrid -- A grid that implementes the Stringer interface
type ASCIIGrid struct {
	*BaseGrid
}

// ASCIIGridHandler -- Defines the interface for a printable grid
type ASCIIGridHandler interface {
	BaseGridHandler
	// ContentsOf Get the content of the cell (what the cell should print inside it)
	ContentsOf(cell *mfp.Cell) string
}

func (g ASCIIGrid) String() string {
	return gridToString(g)
}

func (g ASCIIGrid) ContentsOf(cell *mfp.Cell) string {
	return " "
}
