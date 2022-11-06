package mfp

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"os"
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
	ToImage(cellSize, thickness int, colourTiles bool) *rl.Image
}

// TODO: tile types? use rectangles?

func renderDebugLines(g *Grid, cellSize int) {
	line := rl.Red
	rl.BeginDrawing()
	for x := 0; x < g.columns*cellSize; x += cellSize {
		rl.DrawLine(int32(x), 0, int32(x), int32(g.rows*cellSize), line)
	}
	rl.EndDrawing()
	rl.BeginDrawing()
	for y := 0; y < g.rows*cellSize; y += cellSize {
		rl.DrawLine(0, int32(y), int32(g.columns*cellSize), int32(y), line)
	}
	rl.EndDrawing()
}

func renderWithBezierCurves(g *Grid, cellSize, thickness int) {
	offset := thickness / 2
	wall := rl.Black
	rl.BeginDrawing()
	for cell := range g.EachCell() {
		x1, y1, x2, y2 := (cell.column*cellSize)+offset, (cell.row*cellSize)+offset, ((cell.column+1)*cellSize)+offset, (((cell.row)+1)*cellSize)+offset
		if cell.north == nil {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x2), float32(y1)), float32(thickness), wall)
		}
		if cell.west == nil {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x1), float32(y2)), float32(thickness), wall)
		}
		if !cell.Linked(cell.east) {
			rl.DrawLineBezier(rl.NewVector2(float32(x2), float32(y1)), rl.NewVector2(float32(x2), float32(y2)), float32(thickness), wall)
		}
		if !cell.Linked(cell.south) {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y2)), rl.NewVector2(float32(x2), float32(y2)), float32(thickness), wall)
		}
	}
	rl.EndDrawing()
}

func renderWithTiles(g *Grid, cellSize int) {
	rl.BeginDrawing()
	for cell := range g.EachCell() {
		x, y := int32(cell.column*cellSize), int32(cell.row*cellSize)
		center := rl.NewVector2(float32(x+(int32(cellSize)/2)), float32(y+(int32(cellSize))/2))
		// closed cell
		if !cell.Linked(cell.north) && !cell.Linked(cell.south) && !cell.Linked(cell.east) && !cell.Linked(cell.west) {
			centerX, centerY := x+(int32(cellSize)/2), y+(int32(cellSize))/2
			rl.DrawRectangleLines(x, y, int32(cellSize), int32(cellSize), rl.Black)
			rl.DrawCircleGradient(centerX, centerY, float32(cellSize), rl.Black, rl.White)
			continue
		}
		//--- straight line gradients --
		//north open
		if cell.Linked(cell.north) && !cell.Linked(cell.south) && !cell.Linked(cell.east) && !cell.Linked(cell.west) {
			rl.DrawRectangleGradientV(x, y, int32(cellSize), int32(cellSize), rl.DarkGreen, rl.RayWhite)
			continue
		}
		// South open
		if cell.Linked(cell.south) && !cell.Linked(cell.north) && !cell.Linked(cell.east) && !cell.Linked(cell.west) {
			rl.DrawRectangleGradientV(x, y, int32(cellSize), int32(cellSize), rl.White, rl.Maroon)
			continue
		}
		// East open
		if cell.Linked(cell.east) && !cell.Linked(cell.north) && !cell.Linked(cell.south) && !cell.Linked(cell.west) {
			rl.DrawRectangleGradientH(x, y, int32(cellSize), int32(cellSize), rl.White, rl.Yellow)
			continue
		}
		// west open
		if cell.Linked(cell.west) && !cell.Linked(cell.north) && !cell.Linked(cell.south) && !cell.Linked(cell.east) {
			rl.DrawRectangleGradientH(x, y, int32(cellSize), int32(cellSize), rl.White, rl.DarkGray)
			continue
		}
		//--- diagonal gradients ---
		//We won't actually have any true diagonal gradients
		//since they are way too complex to implement for the scope of this project,
		//so we'll just paint the corners with a slighter
		//reduce alpha value to simulate a gradient
		//North & east open
		if cell.Linked(cell.north) && cell.Linked(cell.east) && !cell.Linked(cell.west) && !cell.Linked(cell.south) {
			fadedColor := rl.Fade(rl.Red, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), fadedColor, rl.Red, fadedColor, rl.White)
			continue
		}
		// east & south open
		if cell.Linked(cell.east) && cell.Linked(cell.south) && !cell.Linked(cell.north) && !cell.Linked(cell.west) {
			fadedColor := rl.Fade(rl.Gold, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), rl.Gold, fadedColor, rl.White, fadedColor)
			continue
		}
		// north & west open
		if cell.Linked(cell.north) && cell.Linked(cell.west) && !cell.Linked(cell.east) && !cell.Linked(cell.south) {
			fadedColor := rl.Fade(rl.Violet, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), rl.White, fadedColor, rl.Violet, fadedColor)
			continue
		}
		// west & south open
		if cell.Linked(cell.west) && cell.Linked(cell.south) && !cell.Linked(cell.north) && !cell.Linked(cell.east) {
			fadedColor := rl.Fade(rl.Pink, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), fadedColor, rl.White, fadedColor, rl.Violet)
			continue
		}
		// --- open ended ---
		// Open Cell
		if cell.Linked(cell.west) && cell.Linked(cell.south) && cell.Linked(cell.north) && cell.Linked(cell.east) {
			rl.DrawCircleGradient(int32(center.X), int32(center.Y), float32(cellSize/2), rl.White, rl.Blue)
			continue
		}
		// -- double openings --
		// We will draw like an hourglass shape superposing two triangles
		// east & west open
		if cell.Linked(cell.east) && cell.Linked(cell.west) && !cell.Linked(cell.north) && !cell.Linked(cell.south) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				rl.NewVector2(float32(x), float32(y+int32(cellSize))),
				center,
				rl.DarkBlue,
			)
			rl.DrawTriangle(
				rl.NewVector2(float32(x+int32(cellSize)), float32(y)),
				center,
				rl.NewVector2(float32(x+int32(cellSize)), float32(y+int32(cellSize))),
				rl.DarkBlue,
			)
			continue
		}
		// north & south open
		if cell.Linked(cell.north) && cell.Linked(cell.south) && !cell.Linked(cell.east) && !cell.Linked(cell.west) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				center,
				rl.NewVector2(float32(x+int32(cellSize)), float32(y)),
				rl.DarkBlue,
			)
			rl.DrawTriangle(
				center,
				rl.NewVector2(
					float32(x),
					float32(y+int32(cellSize)),
				),
				rl.NewVector2(float32(x+int32(cellSize)), float32(y+int32(cellSize))),
				rl.DarkBlue,
			)
			continue
		}
		// --- triple openings ---
		// only north closed
		if !cell.Linked(cell.north) && cell.Linked(cell.south) && cell.Linked(cell.east) && cell.Linked(cell.west) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				center,
				rl.NewVector2(float32(x+int32(cellSize)), float32(y)),
				rl.DarkBlue,
			)
			continue
		}
		// only south closed
		if !cell.Linked(cell.south) && cell.Linked(cell.north) && cell.Linked(cell.east) && cell.Linked(cell.west) {
			rl.DrawTriangle(
				center,
				rl.NewVector2(float32(x), float32(y+int32(cellSize))),
				rl.NewVector2(float32(x+int32(cellSize)), float32(y+int32(cellSize))),
				rl.DarkBlue,
			)
			continue
		}
		// only east closed
		if !cell.Linked(cell.east) && cell.Linked(cell.north) && cell.Linked(cell.south) && cell.Linked(cell.west) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				rl.NewVector2(float32(x), float32(y+int32(cellSize))),
				center,
				rl.DarkBlue,
			)
			continue
		}
		// only west closed
		if !cell.Linked(cell.west) && cell.Linked(cell.north) && cell.Linked(cell.east) && cell.Linked(cell.south) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x+int32(cellSize)), float32(y)),
				center,
				rl.NewVector2(float32(x+int32(cellSize)), float32(y+int32(cellSize))),
				rl.DarkBlue,
			)
			continue
		}
	}
	rl.EndDrawing()
}

func (g *Grid) ToImage(cellSize, thickness int, colourTiles bool) *rl.Image {
	debug := os.Getenv("DEBUG")
	if cellSize <= 0 {
		cellSize = 10
	}
	if thickness <= 0 {
		thickness = 1
	}

	width, height := (cellSize*g.columns)+thickness, (cellSize*g.rows)+thickness
	background := rl.RayWhite

	// Let's use a hidden OpenGL context to
	// draw the image since the texture
	// drawing functions work better
	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(int32(width), int32(height), "")
	rl.SetTargetFPS(60)
	target := rl.LoadRenderTexture(int32(width), int32(height))

	defer rl.EndTextureMode()
	defer rl.UnloadRenderTexture(target)

	rl.BeginTextureMode(target)
	rl.ClearBackground(background)
	if !colourTiles {
		renderWithBezierCurves(g, cellSize, thickness)
	} else {
		renderWithTiles(g, cellSize)
	}
	if debug == "True" {
		renderDebugLines(g, cellSize)
	}

	image := rl.LoadImageFromTexture(target.Texture)
	rl.ImageFlipVertical(*&image)
	return image
}

// There are a few edge cases
// which I cannot quite figure out yet
// since they'll require a cell to
// know the boundaries of the other cells
func (g *Grid) String() string {
	output := "┌" + strings.Repeat("────", g.columns-1) + "───┐" + "\n"
	index := 0
	for row := range g.EachRow() {
		top, bottom := "│", "│"
		for i := range row {
			cell := row[i]
			// three spaces
			body, corner := " c ", "─"
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
