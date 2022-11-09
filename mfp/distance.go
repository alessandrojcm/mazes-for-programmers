package mfp

type Distance struct {
	root  *Cell
	cells map[*Cell]int
}

// DistanceHandler - interface to implement Distance between cells
type DistanceHandler interface {
	PathTo(goal *Cell) Distance
	Max() (maxCell *Cell, maxDistance int)
}

func NewDistance(root *Cell) Distance {
	cells := make(map[*Cell]int)
	cells[root] = 0
	return Distance{
		root:  root,
		cells: cells,
	}
}

// PathTo - Returns the shortest path from the root to the given cell
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

// Max - Returns the cell that is the farthest from the root and its distance value
func (d *Distance) Max() (maxCell *Cell, maxDistance int) {
	maxCell, maxDistance = d.root, 0

	for cell, distance := range d.cells {
		if distance > maxDistance {
			maxCell = cell
			maxDistance = distance
		}
	}
	return
}
