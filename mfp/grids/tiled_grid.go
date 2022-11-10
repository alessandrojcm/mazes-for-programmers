package grids

import rl "github.com/gen2brain/raylib-go/raylib"

type TiledGrid struct {
	*RendererGrid
}

// ToTexture - Returns a texture2d ready to be displayed on the screen or exported to an image, with tiles colored
func (g *TiledGrid) ToTexture(cellSize, thickness int) *rl.RenderTexture2D {
	target := prepareRenderContext(g.columns, g.rows, thickness, cellSize)

	background := rl.RayWhite

	defer rl.EndTextureMode()

	rl.BeginTextureMode(target)
	rl.ClearBackground(background)
	rl.BeginDrawing()
	for cell := range g.EachCell() {
		x, y := int32(cell.Column*cellSize), int32(cell.Row*cellSize)
		center := rl.NewVector2(float32(x+(int32(cellSize)/2)), float32(y+(int32(cellSize))/2))
		// closed cell
		if !cell.Linked(cell.North) && !cell.Linked(cell.South) && !cell.Linked(cell.East) && !cell.Linked(cell.West) {
			centerX, centerY := x+(int32(cellSize)/2), y+(int32(cellSize))/2
			rl.DrawRectangleLines(x, y, int32(cellSize), int32(cellSize), rl.Black)
			rl.DrawCircleGradient(centerX, centerY, float32(cellSize), rl.Black, rl.White)
			continue
		}
		//--- straight line gradients --
		//north open
		if cell.Linked(cell.North) && !cell.Linked(cell.South) && !cell.Linked(cell.East) && !cell.Linked(cell.West) {
			rl.DrawRectangleGradientV(x, y, int32(cellSize), int32(cellSize), rl.DarkGreen, rl.RayWhite)
			continue
		}
		// South open
		if cell.Linked(cell.South) && !cell.Linked(cell.North) && !cell.Linked(cell.East) && !cell.Linked(cell.West) {
			rl.DrawRectangleGradientV(x, y, int32(cellSize), int32(cellSize), rl.White, rl.Maroon)
			continue
		}
		// East open
		if cell.Linked(cell.East) && !cell.Linked(cell.North) && !cell.Linked(cell.South) && !cell.Linked(cell.West) {
			rl.DrawRectangleGradientH(x, y, int32(cellSize), int32(cellSize), rl.White, rl.Yellow)
			continue
		}
		// west open
		if cell.Linked(cell.West) && !cell.Linked(cell.North) && !cell.Linked(cell.South) && !cell.Linked(cell.East) {
			rl.DrawRectangleGradientH(x, y, int32(cellSize), int32(cellSize), rl.White, rl.DarkGray)
			continue
		}
		//--- diagonal gradients ---
		//We won't actually have any true diagonal gradients
		//since they are way too complex to implement for the scope of this project,
		//so we'll just paint the corners with a slighter
		//reduce alpha value to simulate a gradient
		//North & east open
		if cell.Linked(cell.North) && cell.Linked(cell.East) && !cell.Linked(cell.West) && !cell.Linked(cell.South) {
			fadedColor := rl.Fade(rl.Red, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), fadedColor, rl.Red, fadedColor, rl.White)
			continue
		}
		// east & south open
		if cell.Linked(cell.East) && cell.Linked(cell.South) && !cell.Linked(cell.North) && !cell.Linked(cell.West) {
			fadedColor := rl.Fade(rl.Gold, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), rl.Gold, fadedColor, rl.White, fadedColor)
			continue
		}
		// north & west open
		if cell.Linked(cell.North) && cell.Linked(cell.West) && !cell.Linked(cell.East) && !cell.Linked(cell.South) {
			fadedColor := rl.Fade(rl.Violet, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), rl.White, fadedColor, rl.Violet, fadedColor)
			continue
		}
		// west & south open
		if cell.Linked(cell.West) && cell.Linked(cell.South) && !cell.Linked(cell.North) && !cell.Linked(cell.East) {
			fadedColor := rl.Fade(rl.Pink, 0.5)
			rl.DrawRectangleGradientEx(rl.NewRectangle(float32(x), float32(y), float32(cellSize), float32(cellSize)), fadedColor, rl.White, fadedColor, rl.Violet)
			continue
		}
		// --- open ended ---
		// Open Cell
		if cell.Linked(cell.West) && cell.Linked(cell.South) && cell.Linked(cell.North) && cell.Linked(cell.East) {
			rl.DrawCircleGradient(int32(center.X), int32(center.Y), float32(cellSize/2), rl.White, rl.Blue)
			continue
		}
		// -- double openings --
		// We will draw like an hourglass shape superposing two triangles
		// east & west open
		if cell.Linked(cell.East) && cell.Linked(cell.West) && !cell.Linked(cell.North) && !cell.Linked(cell.South) {
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
		if cell.Linked(cell.North) && cell.Linked(cell.South) && !cell.Linked(cell.East) && !cell.Linked(cell.West) {
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
		if !cell.Linked(cell.North) && cell.Linked(cell.South) && cell.Linked(cell.East) && cell.Linked(cell.West) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				center,
				rl.NewVector2(float32(x+int32(cellSize)), float32(y)),
				rl.DarkBlue,
			)
			continue
		}
		// only south closed
		if !cell.Linked(cell.South) && cell.Linked(cell.North) && cell.Linked(cell.East) && cell.Linked(cell.West) {
			rl.DrawTriangle(
				center,
				rl.NewVector2(float32(x), float32(y+int32(cellSize))),
				rl.NewVector2(float32(x+int32(cellSize)), float32(y+int32(cellSize))),
				rl.DarkBlue,
			)
			continue
		}
		// only east closed
		if !cell.Linked(cell.East) && cell.Linked(cell.North) && cell.Linked(cell.South) && cell.Linked(cell.West) {
			rl.DrawTriangle(
				rl.NewVector2(float32(x), float32(y)),
				rl.NewVector2(float32(x), float32(y+int32(cellSize))),
				center,
				rl.DarkBlue,
			)
			continue
		}
		// only west closed
		if !cell.Linked(cell.West) && cell.Linked(cell.North) && cell.Linked(cell.East) && cell.Linked(cell.South) {
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

	return &target
}
