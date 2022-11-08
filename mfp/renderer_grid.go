package mfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

type RendererGrid struct {
	*ASCIIGrid
}

// RendererGridHandler - a grid that can be renderer using raylib
type RendererGridHandler interface {
	ASCIIGridHandler
	ToTexture(cellSize, thickness int, colourTiles bool) *rl.RenderTexture2D
}

func renderDebugLines(g *RendererGrid, cellSize int) {
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

func renderWithBezierCurves(g *RendererGrid, cellSize, thickness int) {
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

func renderWithTiles(g *RendererGrid, cellSize int) {
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

func (g *RendererGrid) ToTexture(cellSize, thickness int, colourTiles bool) *rl.RenderTexture2D {
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

	return &target
}
