package mfp

import "errors"

type Cell struct {
	North, South, East, West *Cell
	Row, Column              int
	links                    map[*Cell]bool
}

type CellHandler interface {
	Link(cell *Cell, bidi bool) error
	Unlink(cell *Cell, bidi bool) error
	Links() []*Cell
	Linked(cell *Cell) bool
	Neighbors() []*Cell
	Distances() Distance
}

func NewCell(row, column int) (*Cell, error) {
	if row < 0 || column < 0 {
		return &Cell{}, errors.New("Row and Column cannot be negative")
	}
	cell := Cell{
		Row:    row,
		Column: column,
	}
	cell.links = make(map[*Cell]bool)
	return &cell, nil
}

func (receiver *Cell) Neighbors() []*Cell {
	// Make with capacity to 4 elements at most
	neighbors := make([]*Cell, 0, 4)

	if receiver.North != nil {
		neighbors = append(neighbors, receiver.North)
	}
	if receiver.East != nil {
		neighbors = append(neighbors, receiver.East)
	}
	if receiver.South != nil {
		neighbors = append(neighbors, receiver.South)
	}
	if receiver.West != nil {
		neighbors = append(neighbors, receiver.West)
	}
	return neighbors
}

func (receiver *Cell) Links() []*Cell {
	var keys []*Cell

	for k := range receiver.links {
		if k == receiver || k == nil {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func (receiver *Cell) Linked(cell *Cell) bool {
	linked, exists := receiver.links[cell]
	// check exists first cause if not !exists
	// then linked is nil
	if !exists {
		return false
	}
	// return linked anyway because it could happen that
	// it exists but is not linked
	return linked
}

func (receiver *Cell) Link(cell *Cell, bidi bool) error {
	receiver.links[cell] = true
	if bidi == true {
		cell.links[receiver] = true
	}
	return nil
}

func (receiver *Cell) Unlink(cell *Cell, bidi bool) error {
	delete(receiver.links, cell)
	if bidi == true {
		return cell.Unlink(receiver, false)
	}
	return nil
}

// Distances A simplified implementation of Dijkstra's algorithm
// TODO: implement bfs
func (receiver *Cell) Distances() Distance {
	distances := NewDistance(receiver)
	frontier := []*Cell{receiver}

	for len(frontier) > 0 {
		var newFrontier []*Cell

		for _, cell := range frontier {
			for _, linked := range cell.Links() {
				_, isLinked := distances.Cells[linked]
				if isLinked {
					continue
				}
				distances.Cells[linked] = distances.Cells[cell] + 1
				newFrontier = append(newFrontier, linked)
			}
		}
		frontier = newFrontier
	}
	return distances
}
