package flags

import rl "github.com/gen2brain/raylib-go/raylib"

func Show(target *rl.RenderTexture2D, name string) {
	rl.ClearWindowState(rl.FlagWindowHidden)
	rl.SetWindowTitle(name)
	defer rl.UnloadTexture(target.Texture)
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.DrawTexture(target.Texture, 0, 0, rl.White)
		rl.EndDrawing()
	}
}
