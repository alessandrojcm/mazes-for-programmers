package mfp

import "errors"

type Cell struct {
	north, south, east, west *Cell
	row, column              int
	links                    map[*Cell]bool
}

type CellHandler interface {
	Link(cell *Cell, bidi bool) error
	Unlink(cell *Cell, bidi bool) error
	Links() []*Cell
	Linked(cell *Cell) bool
	Neighbors() []*Cell
	Distances() Distance[*Cell]
}

func NewCell(row, column int) (*Cell, error) {
	if row < 0 || column < 0 {
		return &Cell{}, errors.New("row and column cannot be negative")
	}
	cell := Cell{
		row:    row,
		column: column,
	}
	cell.links = make(map[*Cell]bool)
	return &cell, nil
}

func (receiver *Cell) Neighbors() []*Cell {
	// Make with capacity to 4 elements at most
	neighbors := make([]*Cell, 0, 4)

	if receiver.north != nil {
		neighbors = append(neighbors, receiver.north)
	}
	if receiver.east != nil {
		neighbors = append(neighbors, receiver.east)
	}
	if receiver.south != nil {
		neighbors = append(neighbors, receiver.south)
	}
	if receiver.west != nil {
		neighbors = append(neighbors, receiver.west)
	}
	return neighbors
}

func (receiver *Cell) Links() []*Cell {
	keys := make([]*Cell, len(receiver.links)-1)

	for k := range receiver.links {
		if k == receiver {
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
func (receiver *Cell) Distances() Distance[*Cell] {
	distances := NewDistance[*Cell](receiver)
	frontier := []*Cell{receiver}

	for len(frontier) > 0 {
		var newFrontier []*Cell

		for _, cell := range frontier {
			for _, linked := range cell.Links() {
				_, isLinked := distances.cells[&linked]
				if isLinked {
					continue
				}
				distances.cells[&linked] = distances.cells[&cell] + 1
				newFrontier = append(newFrontier, linked)
			}
		}
		frontier = newFrontier
	}
	return distances
}
