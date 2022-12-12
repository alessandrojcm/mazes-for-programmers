package algorithms

import (
	"log"
	"math/rand"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"time"
)

// AldousBroder -- an implementation of the Aldous-Broder algorithm
func AldousBroder(grid grids.BaseGridHandler) {
	log.Printf("starting aldous-broder tree run for %dx%d grid", grid.Rows(), grid.Columns())
	defer mfp.TimeTrack(time.Now(), "aldous broder")
	cell, err := grid.RandomCell()
	unvisited := grid.Size() - 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}

	for unvisited > 0 {
		sample := r.Intn(len(cell.Neighbors()))
		log.Println(sample)
		neighbor := cell.Neighbors()[sample]

		if len(neighbor.Links()) == 0 {
			err = cell.Link(neighbor, true)
			if err != nil {
				panic(err)
			}
			unvisited -= 1
			log.Printf("%d cells to go", unvisited)
		}
		cell = neighbor
	}
}
