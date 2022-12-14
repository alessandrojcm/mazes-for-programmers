package mfp

import (
	"errors"
	"log"
	"os"
	"time"
)

// Cell -- basic cell
type Cell struct {
	North, South, East, West *Cell
	Row, Column              int
	links                    map[*Cell]bool
}

type Direction string

const North Direction = "North"
const South Direction = "South"
const East Direction = "East"
const West Direction = "West"

// CellHandler -- This interface defines the bare minimum methods for a basic cell
type CellHandler interface {
	Link(cell *Cell, bidi bool) error
	Unlink(cell *Cell, bidi bool) error
	Links() []*Cell
	Linked(cell *Cell) bool
	Neighbors() []*Cell
	Distances() Distance
	GetDirectionForLink(cell *Cell) (Direction, error)
}

// NewCell -- Factory function for cells
func NewCell(row, column int) (*Cell, error) {
	if row < 0 || column < 0 {
		return &Cell{}, errors.New("row and Column cannot be negative")
	}
	cell := Cell{
		Row:    row,
		Column: column,
	}
	cell.links = make(map[*Cell]bool)
	return &cell, nil
}

// Neighbors -- Get all the surrounding cells
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

// Links -- get all the cells this cell is connected to
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

// Linked -- check if a given cell is connected to this one
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

// Link -- Link this cell to another one
// bidi -- The link is both ways
func (receiver *Cell) Link(cell *Cell, bidi bool) error {
	receiver.links[cell] = true
	if bidi == true {
		cell.links[receiver] = true
	}
	return nil
}

// Unlink -- Unlink this cell from another one
// bidi -- unlink happens both ways
func (receiver *Cell) Unlink(cell *Cell, bidi bool) error {
	delete(receiver.links, cell)
	if bidi == true {
		return cell.Unlink(receiver, false)
	}
	return nil
}

// Distances Gets the weight of the grid using Dijkstra's or BFS (defaults to BFS)
func (receiver *Cell) Distances() Distance {
	// Obscure flag to chose algorithm, for comparison purposes
	alg := os.Getenv("LP_ALG")
	if alg == "dijkstra" {
		return dijkstra(receiver)
	}
	var queue = []*Cell{receiver}
	distances := NewDistance(receiver)
	distances.Cells[receiver] = 0

	log.Printf("Starting weight labelling with bfs")
	defer TimeTrack(time.Now(), "bfs")
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, link := range current.Links() {
			_, isVisited := distances.Cells[link]
			if isVisited {
				continue
			}
			distances.Cells[link] = distances.Cells[current] + 1
			queue = append(queue, link)
		}
	}
	return distances
}

func dijkstra(receiver *Cell) Distance {
	distances := NewDistance(receiver)
	frontier := []*Cell{receiver}
	log.Printf("starting shortest path calculation with Dijkstra for cell %dx%d", receiver.Row, receiver.Column)
	defer TimeTrack(time.Now(), "Dijkstra")
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

// GetDirectionForLink -- this function returns the direction to a given cell that is in this cell's link list
func (receiver *Cell) GetDirectionForLink(cell *Cell) (Direction, error) {
	if receiver.North == cell {
		return North, nil
	}
	if receiver.South == cell {
		return South, nil
	}
	if receiver.East == cell {
		return East, nil
	}
	if receiver.West == cell {
		return West, nil
	}
	return "", errors.New("given cell not in this cell's vicinity")
}
