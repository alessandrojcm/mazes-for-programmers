package mfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"time"
)

// TimeTrack - measure the execution time for a function
// from https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func GetRandomColor() rl.Color {
	return rl.NewColor(
		uint8(rl.GetRandomValue(0, 255)),
		uint8(rl.GetRandomValue(0, 255)),
		uint8(rl.GetRandomValue(0, 255)),
		uint8(255),
	)
}
