package mfp

type Distance struct {
	root  *Cell
	cells map[*Cell]int
}

type DistanceHandler interface {
	PathTo(goal *Cell) Distance
}

func NewDistance(root *Cell) Distance {
	cells := make(map[*Cell]int)
	cells[root] = 0
	return Distance{
		root:  root,
		cells: cells,
	}
}

func (d *Distance) PathTo(goal *Cell) Distance {
	current := goal

	breadcrumbs := NewDistance(d.root)
	breadcrumbs.cells[current] = d.cells[current]

	for current != d.root {
		for _, neighbour := range current.Links() {
			neighbourDistance, isNeighbourLinked := d.cells[neighbour]
			currentDistance, isCurrentLinked := d.cells[current]
			if !isNeighbourLinked || !isCurrentLinked {
				continue
			} else if neighbourDistance < currentDistance {
				breadcrumbs.cells[neighbour] = d.cells[neighbour]
				current = neighbour
				break
			}
		}
	}
	return breadcrumbs
}
