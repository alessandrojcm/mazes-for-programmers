package mfp

import (
	"math/rand"
	"time"
)

const SouthAndWest = "southwest"
const NorthAndWest = "northwest"
const SouthAndEast = "southeast"

func BinaryTree(grid GridHandler, bias string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for cell := range grid.EachCell() {
		neighbors := make([]*Cell, 0, grid.Size())
		switch bias {
		case SouthAndWest:
			if cell.south != nil {
				neighbors = append(neighbors, cell.south)
			}
			if cell.west != nil {
				neighbors = append(neighbors, cell.west)
			}
		case NorthAndWest:
			if cell.north != nil {
				neighbors = append(neighbors, cell.north)
			}
			if cell.west != nil {
				neighbors = append(neighbors, cell.west)
			}
		case SouthAndEast:
			if cell.south != nil {
				neighbors = append(neighbors, cell.south)
			}
			if cell.east != nil {
				neighbors = append(neighbors, cell.east)
			}
		// default bias is northeast
		default:
			if cell.north != nil {
				neighbors = append(neighbors, cell.north)
			}
			if cell.east != nil {
				neighbors = append(neighbors, cell.east)
			}
		}
		randomNeighbor := 0
		if len(neighbors) > 0 {
			randomNeighbor = r.Intn(len(neighbors))
		}
		if len(neighbors) > 0 && neighbors[randomNeighbor] != nil {
			err := cell.Link(neighbors[randomNeighbor], true)
			if err != nil {
				panic(err)
			}
		}
	}
}
