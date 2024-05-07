package grids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"mazes-for-programmers/mfp"
	"os"
	"time"
)

type cellColor struct {
	cell  *mfp.Cell
	color rl.Color
}

func drawMazeLines(lines rl.Texture2D) {
	drawHorizontallyFlipped(lines, rl.Black)
}

// prepareCanvas -- renders a slice of tiles given each tile has a bgColor
func prepareCanvas(cells []cellColor, cellSize, offset int, lines rl.Texture2D, backgroundColor rl.Color) {
	rl.ClearBackground(rl.White)
	rl.BeginDrawing()
	bgColor := backgroundColor
	if bgColor == rl.Blank {
		bgColor = mfp.GetRandomColor()
	}
	// draw every cell first with the solid bg color
	for _, cell := range cells {
		x, y := int32((cell.cell.Column*cellSize)+offset), int32((cell.cell.Row*cellSize)+offset)
		rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), bgColor)
	}
	drawMazeLines(lines)
	rl.EndDrawing()
}

// drawCell -- draws a cell given the cell, the size, the offset and the color
func drawCell(cell *mfp.Cell, cellSize, offset int, color rl.Color) {
	x, y := int32((cell.Column*cellSize)+offset), int32((cell.Row*cellSize)+offset)
	// clear the space first
	rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), rl.White)
	rl.DrawRectangle(x, y, int32(cellSize-offset), int32(cellSize-offset), color)
}

// generateMazeWallsTexture -- draws the lines (walls) for a maze, this returns a texture instead of drawing directly into the target,
// so we can redraw it as necessary
func generateMazeWallsTexture(eachCell chan *mfp.Cell, cellSize, thickness, offset int, width, height int32, wall rl.Color) rl.Texture2D {
	log.Printf("Starting to prerender maze walls")
	defer mfp.TimeTrack(time.Now(), "walls prerendering")
	target := rl.LoadRenderTexture(width, height)
	rl.BeginTextureMode(target)
	rl.BeginDrawing()
	defer rl.EndDrawing()
	defer rl.EndTextureMode()
	for cell := range eachCell {
		x1, y1, x2, y2 := (cell.Column*cellSize)+offset, (cell.Row*cellSize)+offset, ((cell.Column+1)*cellSize)+offset, (((cell.Row)+1)*cellSize)+offset
		if cell.North == nil {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x2), float32(y1)), float32(thickness), wall)
		}
		if cell.West == nil {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x1), float32(y2)), float32(thickness), wall)
		}
		if !cell.Linked(cell.East) {
			rl.DrawLineBezier(rl.NewVector2(float32(x2), float32(y1)), rl.NewVector2(float32(x2), float32(y2)), float32(thickness), wall)
		}
		if !cell.Linked(cell.South) {
			rl.DrawLineBezier(rl.NewVector2(float32(x1), float32(y2)), rl.NewVector2(float32(x2), float32(y2)), float32(thickness), wall)
		}
	}
	return target.Texture
}

// prepareRenderContext -- does some preparations to the opengl context in order to render the maze
func prepareRenderContext(columns, rows, thickness, cellSize int) (target rl.RenderTexture2D) {
	debug := os.Getenv("DEBUG")
	if cellSize <= 0 {
		cellSize = 10
	}
	if thickness <= 0 {
		thickness = 1
	}
	// Let's use a hidden OpenGL context to
	// draw the image since the texture
	// drawing functions work better
	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(0, 0, "")
	rl.SetTargetFPS(30)

	// Making sure the render resolution does not overflow the monitor
	width, height := cellSize*columns, cellSize*rows
	if width > rl.GetMonitorWidth(rl.GetCurrentMonitor()) {
		width = rl.GetMonitorWidth(rl.GetCurrentMonitor())
	}
	if height > rl.GetMonitorHeight(rl.GetCurrentMonitor()) {
		height = rl.GetMonitorHeight(rl.GetCurrentMonitor())
	}
	rl.SetWindowPosition(30, 30)
	rl.SetWindowSize(width, height)

	target = rl.LoadRenderTexture(int32(width), int32(height))

	defer rl.EndTextureMode()
	if debug == "True" {
		line := rl.Red
		rl.BeginDrawing()
		for x := 0; x < columns*cellSize; x += cellSize {
			rl.DrawLine(int32(x), 0, int32(x), int32(rows*cellSize), line)
		}
		rl.EndDrawing()
		rl.BeginDrawing()
		for y := 0; y < rows*cellSize; y += cellSize {
			rl.DrawLine(0, int32(y), int32(columns*cellSize), int32(y), line)
		}
		rl.EndDrawing()
	}
	return
}

// drawHorizontallyFlipped - Flips a texture horizontally
func drawHorizontallyFlipped(texture rl.Texture2D, color rl.Color) {
	// FLIP THE TEXTURE!!
	rl.DrawTextureRec(
		texture,
		rl.NewRectangle(
			0, 0, float32(texture.Width), float32(texture.Height*-1)),
		rl.NewVector2(0, 0),
		color,
	)
}
