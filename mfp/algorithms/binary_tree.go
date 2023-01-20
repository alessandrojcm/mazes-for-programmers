package algorithms

import (
	"log"
	"math/rand"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

const SouthAndWest = "southwest"
const NorthAndWest = "northwest"
const SouthAndEast = "southeast"

// BinaryTree -- a binary tree maze implementation, with bias to change the texture of the maze
func BinaryTree(grid grids.BaseGridHandler, bias string, cutOffPoint int) {
	defer mfp.TimeTrack(time.Now(), "binary tree run")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	log.Printf("starting binary tree run for %dx%d grid", grid.Rows(), grid.Columns())
	count := 0
	for cell := range grid.EachCell() {
		if cutOffPoint != -1 && count >= cutOffPoint {
			return
		}
		neighbors := make([]*mfp.Cell, 0, grid.Size())
		switch bias {
		case SouthAndWest:
			if cell.South != nil {
				neighbors = append(neighbors, cell.South)
			}
			if cell.West != nil {
				neighbors = append(neighbors, cell.West)
			}
		case NorthAndWest:
			if cell.North != nil {
				neighbors = append(neighbors, cell.North)
			}
			if cell.West != nil {
				neighbors = append(neighbors, cell.West)
			}
		case SouthAndEast:
			if cell.South != nil {
				neighbors = append(neighbors, cell.South)
			}
			if cell.East != nil {
				neighbors = append(neighbors, cell.East)
			}
		// default bias is northeast
		default:
			if cell.North != nil {
				neighbors = append(neighbors, cell.North)
			}
			if cell.East != nil {
				neighbors = append(neighbors, cell.East)
			}
		}
		randomNeighbor := 0
		if len(neighbors) > 0 {
			randomNeighbor = r.Intn(len(neighbors))
		}
		if len(neighbors) > 0 && neighbors[randomNeighbor] != nil {
			count++
			err := cell.Link(neighbors[randomNeighbor], true)
			if err != nil {
				panic(err)
			}
		}
	}
}
