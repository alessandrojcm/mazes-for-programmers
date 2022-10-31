package mfp

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"strings"
	"time"
)

type Grid struct {
	rows, columns int
	cells         [][]*Cell
}

type GridHandler interface {
	PrepareGrid() error
	ConfigureCells() error
	CellAt(row, column int) (*Cell, error)
	Size() int
	RandomCell() (*Cell, error)
	EachRow() chan []*Cell
	EachCell() chan *Cell
	ToPng(cellSize int) error
}

func (g *Grid) ToPng(cellSize int) error {
	if cellSize <= 0 {
		cellSize = 10
	}
	width, height := cellSize*g.columns, cellSize*g.rows
	background, wall := rl.White, rl.Black
	image := rl.GenImageColor(width, height, background)
	// Draw a "frame" around the image
	rl.ImageDrawRectangleLines(image, rl.NewRectangle(0, 0, float32(width), float32(height)), 1, wall)
	for cell := range g.EachCell() {
		x1, y1, x2, y2 := cell.column*cellSize, cell.row*cellSize, (cell.column+1)*cellSize, ((cell.row)+1)*cellSize
		if cell.north == nil {
			rl.ImageDrawLine(image, int32(x1), int32(y1), int32(x2), int32(y1), wall)
		}
		if cell.west == nil {
			rl.ImageDrawLine(image, int32(x1), int32(y1), int32(x1), int32(y2), wall)
		}
		if !cell.Linked(cell.east) {
			rl.ImageDrawLine(image, int32(x2), int32(y1), int32(x2), int32(y2), wall)
		}
		if !cell.Linked(cell.south) {
			rl.ImageDrawLine(image, int32(x1), int32(y2), int32(x2), int32(y2), wall)
		}
	}
	rl.ImageFlipVertical(*&image)
	rl.ExportImage(*image, "maze.png")
	rl.UnloadImage(image)
	return nil
}

func (g *Grid) String() string {
	output := "+" + strings.Repeat("----", g.columns-1) + "---+" + "\n"
	for row := range g.EachRow() {
		top, bottom := "|", "b"
		for i := range row {
			cell := row[i]
			// three spaces
			body, corner := " c ", "|"
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
				eastBoundary = "|"
			}
			top = top + body + eastBoundary
			if cell.Linked(cell.south) {
				southBoundary = "   "
			} else {
				southBoundary = "---"
			}
			if !cell.Linked(cell.east) && !cell.Linked(cell.south) {
				corner = "+"
			}
			if cell.Linked(cell.south) && !cell.Linked(cell.east) {
				corner = "|"
			}
			if cell.Linked(cell.east) && cell.Linked(cell.south) {
				corner = "-"
			}
			if cell.Linked(cell.east) && !cell.Linked(cell.south) {
				corner = "-"
			}
			bottom = bottom + southBoundary + corner
		}
		output = output + top + "\n" + bottom + "\n"
	}
	return output
}

func NewGrid(rows, columns int) (*Grid, error) {
	if rows < 0 || columns < 0 {
		return &Grid{}, errors.New("rows and columns must be greater than 0")

	}
	newGrid := Grid{
		rows:    rows,
		columns: columns,
	}
	err := newGrid.PrepareGrid()
	if err != nil {
		return &Grid{}, err
	}
	err = newGrid.ConfigureCells()
	if err != nil {
		return &Grid{}, err
	}
	return &newGrid, nil
}

func (g *Grid) CellAt(row, column int) (*Cell, error) {
	if row < 0 || row > g.rows-1 || column < 0 || column > g.columns-1 {
		return &Cell{}, errors.New("Cell out of bounds")
	}
	return g.cells[row][column], nil
}

func (g *Grid) PrepareGrid() error {
	g.cells = make([][]*Cell, g.rows)
	for row := range g.cells {
		g.cells[row] = make([]*Cell, g.columns)
		for column := range g.cells[row] {
			cell, err := NewCell(row, column)
			if err != nil {
				return err
			}
			g.cells[row][column] = cell
		}
	}
	return nil
}

func (g *Grid) ConfigureCells() error {
	for _, row := range g.cells {
		for _, cell := range row {
			rowNum, colNum := cell.row, cell.column
			// we need to do boundary checking,
			// there is no operator overriding
			// in Go so let's do it the old-fashioned way
			// if any calculation goes out of bounds
			// just assign nil
			if rowNum-1 < 0 {
				cell.north = nil
			} else {
				cell.north = g.cells[rowNum-1][colNum]
			}
			if rowNum+1 > g.rows-1 {
				cell.south = nil
			} else {
				cell.south = g.cells[rowNum+1][colNum]
			}
			if colNum-1 < 0 {
				cell.west = nil
			} else {
				cell.west = g.cells[rowNum][colNum-1]
			}
			if colNum+1 > g.columns-1 {
				cell.east = nil
			} else {
				cell.east = g.cells[rowNum][colNum+1]
			}
		}
	}
	return nil
}

func (g *Grid) RandomCell() (*Cell, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	row, column := r.Intn(g.rows), r.Intn(g.columns)

	return g.cells[row][column], nil
}

func (g *Grid) Size() int {
	return g.rows * g.columns
}

func (g *Grid) EachRow() chan []*Cell {
	rowChan := make(chan []*Cell)
	go func() {
		for _, row := range g.cells {
			rowChan <- row
		}
		close(rowChan)
	}()
	return rowChan
}

func (g *Grid) EachCell() chan *Cell {
	cellChan := make(chan *Cell)
	go func() {
		for row := 0; row < g.rows; row++ {
			for column := 0; column < g.columns; column++ {
				cell, err := g.CellAt(row, column)
				if err != nil {
					panic(err)
				}
				cellChan <- cell
			}
		}
		close(cellChan)
	}()

	return cellChan
}
