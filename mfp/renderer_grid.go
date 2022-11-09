package mfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

type RendererGrid struct {
	*ASCIIGrid
}

// RendererGridHandler - a grid that can be rendered using raylib, the default render method is using lines
type RendererGridHandler interface {
	ASCIIGridHandler
	ToTexture(cellSize, thickness int) *rl.RenderTexture2D
	BackgroundColorForCell(cell *Cell) color.RGBA
}

func (g *RendererGrid) BackgroundColorForCell(cell *Cell) rl.Color {
	return rl.Blank
}

// ToTexture - Returns a texture2d ready to be displayed on the screen or exported to an image rendered with lines and walls
func (g *RendererGrid) ToTexture(cellSize, thickness int) *rl.RenderTexture2D {
	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)

	background := rl.RayWhite

	defer rl.EndTextureMode()

	rl.BeginTextureMode(target)
	rl.ClearBackground(background)
	offset := thickness / 2
	wall := rl.Black
	drawMazeLines(g.EachCell(), cellSize, thickness, offset, wall)

	return &target
}
