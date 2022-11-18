package grids

import (
	"fmt"
	"log"
	"mazes-for-programmers/mfp"
	"strings"
	"time"
)

// There are a few edge cases
// which I cannot quite figure out yet
// since they'll require a cell to
// know the boundaries of the other cells
func gridToString(g ASCIIGridHandler) string {
	output := "┌" + strings.Repeat("────", g.Columns()-1) + "───┐" + "\n"
	index := 0
	log.Printf("starting to print grid with %dx%d dimention", g.Rows(), g.Columns())
	defer mfp.TimeTrack(time.Now(), "grid printing")
	for row := range g.EachRow() {
		top, bottom := "│", "│"
		for i := range row {
			cell := row[i]
			// three spaces
			body, corner := fmt.Sprintf(" %s ", g.ContentsOf(cell)), "─"
			var eastBoundary, southBoundary string

			if cell == nil {
				cell = &mfp.Cell{
					Row:    -1,
					Column: -1,
				}
			}
			if cell.Linked(cell.East) {
				eastBoundary = " "
			} else {
				eastBoundary = "│"
			}
			top = top + body + eastBoundary
			if cell.Linked(cell.South) {
				southBoundary = "   "
			} else {
				southBoundary = "───"
			}
			// special case so just short-circuit the loop
			if index == g.Rows()-1 && i == 0 {
				bottom = "└" + southBoundary + corner
				continue
			} else if index == g.Rows()-1 && i == len(row)-1 {
				bottom = bottom + southBoundary + "┘"
				continue
			}
			if cell.Linked(cell.South) && cell.Linked(cell.East) {
				corner = "╷"
			}
			if !cell.Linked(cell.South) && cell.Linked(cell.East) {
				corner = "─"
			}
			if cell.Linked(cell.South) && !cell.Linked(cell.East) {
				if index == g.Rows()-1 {
					corner = "─"
				} else if i == len(row)-1 {
					corner = "│"
				} else {
					corner = "╷"
				}
			}
			if !cell.Linked(cell.South) && !cell.Linked(cell.East) {
				if index == g.Rows()-1 {
					corner = "┴"
				} else {
					corner = "│"
				}
			}

			bottom = bottom + southBoundary + corner
		}
		index++
		output = output + top + "\n" + bottom + "\n"
	}
	return output
}
