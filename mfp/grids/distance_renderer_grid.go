package grids

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"log"
	"mazes-for-programmers/mfp"
	"time"
)

// DistanceRenderGrid -- a grid that can render the distance weight using raylib
type DistanceRenderGrid struct {
	backgroundColor rl.Color
	*DistanceGrid
}

type DistanceRendererGridHandler interface {
	*RendererGridHandler
}

// BackgroundColorForCell -- Computes an alpha value taking the distance into account and paints the background
// accordingly (more translucent means farther)
func (g *DistanceRenderGrid) BackgroundColorForCell(cell *mfp.Cell) color.RGBA {
	distance, isOk := g.Distances.Cells[cell]
	bgColor := g.backgroundColor
	// blank means random color
	// so we'll generate one
	if bgColor == rl.Blank {
		bgColor = mfp.GetRandomColor()
	}
	if !isOk {
		return bgColor
	}
	_, maximum := g.Distances.Max()
	intensity := (float32(maximum) - float32(distance)) / float32(maximum)
	return rl.Fade(bgColor, intensity)
}

func (g *DistanceRenderGrid) ToTexture(cellSize, thickness int, printWeights bool) *rl.RenderTexture2D {
	log.Printf("starting to render distenced grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "grid rendering")
	offset := thickness / 2
	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)
	lines := generateMazeWallsTexture(g.EachCell(), cellSize, thickness, offset, target.Texture.Width, target.Texture.Height, rl.Black)

	rl.BeginTextureMode(target)
	defer rl.EndDrawing()
	defer rl.EndTextureMode()

	var cells []struct {
		x      int32
		y      int32
		weight int
		color  rl.Color
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for cell := range g.DistanceGrid.EachCell() {
		x, y := int32((cell.Column*cellSize)+offset), int32((cell.Row*cellSize)+offset)
		c := g.BackgroundColorForCell(cell)
		rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), c)
		if printWeights {
			cells = append(cells, struct {
				x      int32
				y      int32
				weight int
				color  rl.Color
			}{x: x, y: y, color: c, weight: g.Distances.Cells[cell]})
		}
	}

	drawHorizontallyFlipped(lines, rl.Black)

	rl.EndDrawing()
	rl.EndTextureMode()
	if !printWeights {
		return &target
	}
	// The main texture is horizontally flipped since OpenGL takes the bottom left corner as the origin.
	// Hence, if we draw the text at once it will also be inverted.
	// To get aroung this, we draw the lines FIRST and then we draw the text, un flipped, in individual
	// textures that will be added on top of the main texture.
	for _, c := range cells {
		if c.weight == 0 {
			continue
		}
		dimension := int32(cellSize - 10)
		weightCell := rl.LoadRenderTexture(dimension, dimension)
		rl.BeginTextureMode(weightCell)
		rl.BeginDrawing()
		rl.DrawText(fmt.Sprintf("%d", c.weight), int32(offset), int32(offset), int32(cellSize/2), rl.Black)
		rl.EndDrawing()
		rl.EndTextureMode()

		rl.BeginTextureMode(target)
		rl.BeginDrawing()
		drawHorizontallyFlipped(weightCell.Texture, rl.White)
		rl.EndDrawing()
		rl.EndTextureMode()
	}
	return &target
}
