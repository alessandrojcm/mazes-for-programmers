package grids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"mazes-for-programmers/mfp"
	"sort"
	"time"
)

type BreadcrumbGrid struct {
	*AnimatableGrid
}

// drawPath -- draws a straight line given a cell and two different points
func drawPath(cell *mfp.Cell, cellSize, offset int, from, to mfp.Direction) {
	x, y := float32((cell.Column*cellSize)+offset), float32((cell.Row*cellSize)+offset)
	halfStep := float32(cellSize-offset) / 2
	getVectorPath := func(dir mfp.Direction) rl.Vector2 {
		switch dir {
		case mfp.North:
			return rl.NewVector2(
				x+halfStep,
				y,
			)
		case mfp.East:
			return rl.NewVector2(
				x+halfStep*2,
				y+halfStep,
			)
		case mfp.South:
			return rl.NewVector2(
				x+halfStep,
				y+halfStep*2,
			)
		default:
			return rl.NewVector2(x, y+halfStep)
		}
	}
	rl.DrawLineV(getVectorPath(from), getVectorPath(to), rl.Red)
}

func (g *BreadcrumbGrid) ShowAnimation(cellSize, thickness int) {
	log.Printf("starting to animate breadcrumb grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "breadcrumb grid rendering")

	var cells []cellColor
	var breadCrumbCells []cellColor
	isFinished := false
	hint := "Press R to restart"

	// We will follow the same logic as the animatable grid, using the transparency to
	// get the cells for the path; the difference is that the
	// path will be used to draw the line
	for cell := range g.DistanceGrid.EachCell() {
		c := cellColor{
			cell:  cell,
			color: g.BackgroundColorForCell(cell),
		}
		cells = append(cells, c)
		if c.color.A < uint8(255) {
			breadCrumbCells = append(breadCrumbCells, c)
		}
	}
	// sort the cells to draw the line
	sort.Slice(breadCrumbCells, func(i, j int) bool {
		return breadCrumbCells[i].color.A < breadCrumbCells[j].color.A
	})

	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)
	offset := thickness / 2
	// show the window
	rl.ClearWindowState(rl.FlagWindowHidden)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	lines := generateMazeWallsTexture(g.EachCell(), cellSize, thickness, offset, target.Texture.Width, target.Texture.Height, rl.Black)
	timer := float32(0.0)

	// keep the cell we're currently drawing
	actualCell := 0
	prepareCanvas(cells, cellSize, offset, lines, g.backgroundColor)
	// now animate the cells per frame
	for !rl.WindowShouldClose() {
		if isFinished && rl.IsKeyPressed(rl.KeyR) {
			prepareCanvas(cells, cellSize, offset, lines, g.backgroundColor)
			actualCell = 0
			isFinished = false
		}
		// We've finished, don't do anything
		if isFinished {
			rl.BeginDrawing()
			rl.DrawText(hint, 5, target.Texture.Height-29, 24, rl.Black)
			rl.EndDrawing()
			continue
		}
		timer += rl.GetFrameTime()
		if timer < 0.05 {
			continue
		}
		timer = 0.0
		rl.BeginDrawing()
		if actualCell > 0 && actualCell < len(breadCrumbCells)-2 {
			tail, _ := breadCrumbCells[actualCell].cell.GetDirectionForLink(breadCrumbCells[actualCell-1].cell)
			middle, _ := breadCrumbCells[actualCell].cell.GetDirectionForLink(breadCrumbCells[actualCell+1].cell)
			drawPath(breadCrumbCells[actualCell].cell, cellSize, offset, tail, middle)
			log.Printf("cell %d from the %s to the %s\n", actualCell, tail, middle)
		}
		actualCell++
		rl.EndDrawing()
		isFinished = actualCell > actualCell%len(breadCrumbCells)
	}
	rl.EndTextureMode()
}
