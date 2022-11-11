package grids

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"log"
	"mazes-for-programmers/mfp"
	"time"
)

type DistanceRenderGrid struct {
	*DistanceGrid
}

func (g *DistanceRenderGrid) BackgroundColorForCell(cell *mfp.Cell) color.RGBA {
	distance, isOk := g.Distances.Cells[cell]
	if !isOk {
		return rl.DarkGreen
	}
	_, maximum := g.Distances.Max()
	intensity := (float32(maximum) - float32(distance)) / float32(maximum)
	fmt.Println(intensity)
	return rl.Fade(rl.DarkGreen, intensity)
}

func (g *DistanceRenderGrid) ToTexture(cellSize, thickness int) *rl.RenderTexture2D {
	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)
	offset := thickness / 2

	rl.BeginTextureMode(target)
	defer rl.EndTextureMode()
	log.Printf("starting to render distenced grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "grid rendering")
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for cell := range g.DistanceGrid.EachCell() {
		x, y := int32((cell.Column*cellSize)+offset), int32((cell.Row*cellSize)+offset)
		rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), g.BackgroundColorForCell(cell))
	}
	rl.EndDrawing()
	drawMazeLines(g.EachCell(), cellSize, thickness, offset, rl.Black)
	return &target
}
