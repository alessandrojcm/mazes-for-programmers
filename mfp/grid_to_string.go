package mfp

import (
	"fmt"
	"strings"
)

func gridToString(g ASCIIGridHandler) string {
	output := "┌" + strings.Repeat("────", g.Columns()-1) + "───┐" + "\n"
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
			if index == g.Rows()-1 && i == 0 {
				bottom = "└" + southBoundary + corner
				continue
			} else if index == g.Rows()-1 && i == len(row)-1 {
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
				if index == g.Rows()-1 {
					corner = "─"
				} else if i == len(row)-1 {
					corner = "│"
				} else {
					corner = "╷"
				}
			}
			if !cell.Linked(cell.south) && !cell.Linked(cell.east) {
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
