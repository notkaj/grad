// Package animation provides animations (duh!)
package animation

import (
	"math/rand"
	"time"

	"charm.land/bubbles/v2/spinner"
)

func Playing() spinner.Spinner {
	return spinner.Spinner{
		Frames: frames(),
		FPS:    time.Second / 7,
	}
}

func frames() []string {
	var frames [20]string
	for f := range frames {
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
