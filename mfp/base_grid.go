package mfp

import (
	"errors"
	"math/rand"
	"time"
)

type BaseGrid struct {
	rows, columns int
	cells         [][]*Cell
}

type BaseGridHandler interface {
	PrepareGrid() error
	ConfigureCells() error
	CellAt(row, column int) (*Cell, error)
	Size() int
	RandomCell() (*Cell, error)
	EachRow() chan []*Cell
	EachCell() chan *Cell
	Empty() bool
	Rows() int
	Columns() int
}

func (g *BaseGrid) Rows() int {
	return g.rows
}

func (g *BaseGrid) Columns() int {
	return g.columns
}

func (g *BaseGrid) Empty() bool {
	return g.rows == 0 || g.columns == 0 || len(g.cells) == 0
}

func (g *BaseGrid) CellAt(row, column int) (*Cell, error) {
	if row < 0 || row > g.rows-1 || column < 0 || column > g.columns-1 {
		return &Cell{}, errors.New("Cell out of bounds")
	}
	return g.cells[row][column], nil
}

func (g *BaseGrid) PrepareGrid() error {
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

func (g *BaseGrid) ConfigureCells() error {
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

func (g *BaseGrid) RandomCell() (*Cell, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	row, column := r.Intn(g.rows), r.Intn(g.columns)

	return g.cells[row][column], nil
}

func (g *BaseGrid) Size() int {
	return g.rows * g.columns
}

func (g *BaseGrid) EachRow() chan []*Cell {
	rowChan := make(chan []*Cell)
	go func() {
		for _, row := range g.cells {
			rowChan <- row
		}
		close(rowChan)
	}()
	return rowChan
}

func (g *BaseGrid) EachCell() chan *Cell {
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
