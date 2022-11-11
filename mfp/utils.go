package mfp

import (
	"log"
	"time"
)

// TimeTrack - measure the execution time for a function
// from https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
