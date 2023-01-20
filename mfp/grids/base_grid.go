package grids

import (
	"errors"
	"github.com/gookit/goutil/arrutil"
	"log"
	"mazes-for-programmers/mfp"
	"time"
)

// BaseGrid -- A basic grid with no means to print or otherwise visually show itself
type BaseGrid struct {
	rows, columns int
	cells         [][]*mfp.Cell
}

// BaseGridHandler -- defines the very basic methods for any grid
type BaseGridHandler interface {
	PrepareGrid() error
	ConfigureCells() error
	CellAt(row, column int) (*mfp.Cell, error)
	Size() int
	RandomCell() (*mfp.Cell, error)
	EachRow() chan []*mfp.Cell
	EachCell() chan *mfp.Cell
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

// CellAt -- accessor to prevent out of bound errors
func (g *BaseGrid) CellAt(row, column int) (*mfp.Cell, error) {
	if row < 0 || row > g.rows-1 || column < 0 || column > g.columns-1 {
		return &mfp.Cell{}, errors.New("cell out of bounds")
	}
	return g.cells[row][column], nil
}

// PrepareGrid -- initializes the grid, filling the cells
func (g *BaseGrid) PrepareGrid() error {
	g.cells = make([][]*mfp.Cell, g.rows)
	log.Printf("starting to prepare grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "grid preparation")
	for row := range g.cells {
		g.cells[row] = make([]*mfp.Cell, g.columns)
		for column := range g.cells[row] {
			cell, err := mfp.NewCell(row, column)
			if err != nil {
				return err
			}
			g.cells[row][column] = cell
		}
	}
	return nil
}

// ConfigureCells -- Sets the links for the cells (which cells do a specific cell has a boundary with)
func (g *BaseGrid) ConfigureCells() error {
	log.Printf("starting to configure grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "grid configuration")
	for _, row := range g.cells {
		for _, cell := range row {
			rowNum, colNum := cell.Row, cell.Column
			// we need to do boundary checking,
			// there is no operator overriding
			// in Go so let's do it the old-fashioned way
			// if any calculation goes out of bounds
			// just assign nil
			if rowNum-1 < 0 {
				cell.North = nil
			} else {
				cell.North = g.cells[rowNum-1][colNum]
			}
			if rowNum+1 > g.rows-1 {
				cell.South = nil
			} else {
				cell.South = g.cells[rowNum+1][colNum]
			}
			if colNum-1 < 0 {
				cell.West = nil
			} else {
				cell.West = g.cells[rowNum][colNum-1]
			}
			if colNum+1 > g.columns-1 {
				cell.East = nil
			} else {
				cell.East = g.cells[rowNum][colNum+1]
			}
		}
	}
	return nil
}

func (g *BaseGrid) RandomCell() (*mfp.Cell, error) {
	return arrutil.GetRandomOne(arrutil.GetRandomOne[[]*mfp.Cell](g.cells)), nil
}

func (g *BaseGrid) Size() int {
	return g.rows * g.columns
}

func (g *BaseGrid) EachRow() chan []*mfp.Cell {
	rowChan := make(chan []*mfp.Cell)
	go func() {
		for _, row := range g.cells {
			rowChan <- row
		}
		close(rowChan)
	}()
	return rowChan
}

func (g *BaseGrid) EachCell() chan *mfp.Cell {
	cellChan := make(chan *mfp.Cell)
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
