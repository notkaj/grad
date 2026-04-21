// Package animation provides animations (duh!)
package animation

import (
	"math/rand"
	"time"

	"charm.land/bubbles/v2/spinner"
)

var PlayingSpinner = spinner.Spinner{
	Frames: frames(),
	FPS:    time.Second / 10,
}

func frames() []string {
	var frames [5]string
	for f := range 5 {
		var runes [3]rune
		for r := range 3 {
			i := rand.Intn(8)
			runes[r] = b[i]
		}
		frames[f] = string(runes[:])
	}
	return frames[:]
}

var b = []rune{
	'▁',
	'▂',
	'▃',
	'▄',
	'▅',
	'▆',
	'▇',
	'█',
}
