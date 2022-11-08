package mfp

type Distance struct {
	root  *Cell
	cells map[*Cell]int
}

func NewDistance(root *Cell) Distance {
	cells := make(map[*Cell]int)
	cells[root] = 0
	return Distance{
		root:  root,
		cells: cells,
	}
}
