package grids

import (
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

func (g *DistanceRenderGrid) ToTexture(cellSize, thickness int) *rl.RenderTexture2D {
	log.Printf("starting to render distenced grid with %dx%d dimention", g.rows, g.columns)
	defer mfp.TimeTrack(time.Now(), "grid rendering")
	offset := thickness / 2
	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)
	lines := generateMazeWallsTexture(g.EachCell(), cellSize, thickness, offset, target.Texture.Width, target.Texture.Height, rl.Black)

	rl.BeginTextureMode(target)
	defer rl.EndDrawing()
	defer rl.EndTextureMode()

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for cell := range g.DistanceGrid.EachCell() {
		x, y := int32((cell.Column*cellSize)+offset), int32((cell.Row*cellSize)+offset)
		rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), g.BackgroundColorForCell(cell))
	}

	// FLIP THE TEXTURE!!
	rl.DrawTextureRec(
		lines,
		rl.NewRectangle(
			0, 0, float32(lines.Width), float32(lines.Height*-1)),
		rl.NewVector2(0, 0),
		rl.Black,
	)
	return &target
}
