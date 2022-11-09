package mfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

func drawMazeLines(eachCell chan *Cell, cellSize, thickness, offset int, wall rl.Color) {
	rl.BeginDrawing()
	for cell := range eachCell {
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

func prepareRenderContext(columns, rows, thickness, cellSize int) (target rl.RenderTexture2D) {
	debug := os.Getenv("DEBUG")
	if cellSize <= 0 {
		cellSize = 10
	}
	if thickness <= 0 {
		thickness = 1
	}

	width, height := (cellSize*columns)+thickness, (cellSize*rows)+thickness

	// Let's use a hidden OpenGL context to
	// draw the image since the texture
	// drawing functions work better
	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(int32(width), int32(height), "")
	rl.SetTargetFPS(60)
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
