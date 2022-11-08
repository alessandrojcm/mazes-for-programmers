package mfp

import (
	"fmt"
	"strings"
)

type ASCIIGrid struct {
	*BaseGrid
}

type ASCIIGridHandler interface {
	BaseGridHandler
	ContentsOf(cell *Cell) string
}

// There are a few edge cases
// which I cannot quite figure out yet
// since they'll require a cell to
// know the boundaries of the other cells
func (g *ASCIIGrid) String() string {
	output := "┌" + strings.Repeat("────", g.columns-1) + "───┐" + "\n"
	index := 0
	for row := range g.EachRow() {
		top, bottom := "│", "│"
		for i := range row {
			cell := row[i]
			// three spaces
			body, corner := fmt.Sprintf(" %s ", g.ContentsOf(cell)), "─"
			var eastBoundary, southBoundary string

			if cell == nil {
				cell = &Cell{
					row:    -1,
					column: -1,
				}
			}
			if cell.Linked(cell.east) {
				eastBoundary = " "
			} else {
				eastBoundary = "│"
			}
			top = top + body + eastBoundary
			if cell.Linked(cell.south) {
				southBoundary = "   "
			} else {
				southBoundary = "───"
			}
			// special case so just short-circuit the loop
			if index == g.rows-1 && i == 0 {
				bottom = "└" + southBoundary + corner
				continue
			} else if index == g.rows-1 && i == len(row)-1 {
				bottom = bottom + southBoundary + "┘"
				continue
			}
			if cell.Linked(cell.south) && cell.Linked(cell.east) {
				corner = "╷"
			}
			if !cell.Linked(cell.south) && cell.Linked(cell.east) {
				corner = "─"
			}
			if cell.Linked(cell.south) && !cell.Linked(cell.east) {
				if index == g.rows-1 {
					corner = "─"
				} else if i == len(row)-1 {
					corner = "│"
				} else {
					corner = "╷"
				}
			}
			if !cell.Linked(cell.south) && !cell.Linked(cell.east) {
				if index == g.rows-1 {
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

func (g *ASCIIGrid) ContentsOf(cell *Cell) string {
	return " "
}
