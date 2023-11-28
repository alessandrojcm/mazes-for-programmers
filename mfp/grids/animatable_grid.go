package grids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"mazes-for-programmers/mfp"
	"sort"
	"time"
)

type AnimatableGrid struct {
	*DistanceRenderGrid
}

type AnimatableGridHandler interface {
	RendererGridHandler
	ShowAnimation(cellSize, thickness int)
}

func (g *AnimatableGrid) ShowAnimation(cellSize, thickness int) {
	log.Printf("starting to animate distanced grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "animatable grid rendering")

	var cells []cellColor
	var translucentCells []cellColor
	isFinished := false
	hint := "Press R to restart"

	// let's precompute all the cell's colours, and set aside the translucent ones in another array
	for cell := range g.DistanceGrid.EachCell() {
		c := cellColor{
			cell:  cell,
			color: g.BackgroundColorForCell(cell),
		}
		cells = append(cells, c)
		if c.color.A < uint8(255) {
			translucentCells = append(translucentCells, c)
		}
	}
	sort.Slice(translucentCells, func(i, j int) bool {
		return translucentCells[i].color.A < translucentCells[j].color.A
	})

	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)
	offset := thickness / 2
	// show the window
	rl.ClearWindowState(rl.FlagWindowHidden)

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
		drawCell(translucentCells[actualCell].cell, cellSize, offset, translucentCells[actualCell].color)
		actualCell++
		drawMazeLines(lines)
		rl.EndDrawing()
		isFinished = actualCell > actualCell%len(translucentCells)
	}
	rl.EndTextureMode()
}
