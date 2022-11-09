package mfp

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

type DistanceRenderGrid struct {
	*DistanceGrid
}

func (g *DistanceRenderGrid) BackgroundColorForCell(cell *Cell) color.RGBA {
	distance, isOk := g.Distances.cells[cell]
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
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for cell := range g.DistanceGrid.EachCell() {
		// reducing the overall surface of the square since they get painted
		// after the lines, otherwise they would get painted on top of the lines
		x, y := int32((cell.column*cellSize)+offset), int32((cell.row*cellSize)+offset)
		rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), g.BackgroundColorForCell(cell))
	}
	rl.EndDrawing()
	drawMazeLines(g.EachCell(), cellSize, thickness, offset, rl.Black)
	return &target
}
