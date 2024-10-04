package sleep

import (
	"math/rand"
	"testing"
	"time"
)

func Rand(d time.Duration) {
	if testing.Testing() {
		return
	}

	randomFactor := 0.1 + rand.Float64()
	randomDuration := time.Duration(float64(d) * randomFactor)
	time.Sleep(randomDuration)
}

func Sleep(d time.Duration) {
	if testing.Testing() {
		return
	}

	time.Sleep(d)
}
