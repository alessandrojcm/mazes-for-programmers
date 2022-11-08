package mfp

type Constraint = CellHandler

type Distance[TCell Constraint] struct {
	root  *TCell
	cells map[*TCell]int
}

func NewDistance[TCell Constraint](root TCell) Distance[TCell] {
	cells := make(map[*TCell]int)
	cells[&root] = 0
	return Distance[TCell]{
		root:  &root,
		cells: cells,
	}
}
