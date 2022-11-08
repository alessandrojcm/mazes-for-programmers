package mfp

type ASCIIGrid struct {
	*BaseGrid
}

type ASCIIGridHandler interface {
	BaseGridHandler
	ContentsOf(cell *Cell) string
}

// There are a few edge cases
// which I cannot quite figure out yet
// since they'll require a cell to
// know the boundaries of the other cells
func (g ASCIIGrid) String() string {
	return gridToString(g)
}

func (g ASCIIGrid) ContentsOf(cell *Cell) string {
	return " "
}
