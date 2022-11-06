package flags

import rl "github.com/gen2brain/raylib-go/raylib"

func Export(target *rl.RenderTexture2D, name string) {
	img := rl.LoadImageFromTexture(target.Texture)
	rl.ImageFlipVertical(*&img)
	defer rl.UnloadImage(img)
	rl.ExportImage(*img, name)
}
