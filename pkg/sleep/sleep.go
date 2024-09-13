package sleep

import (
	"math/rand"
	"time"
)

func Rand(d time.Duration) {
	randomFactor := 0.1 + rand.Float64()
	randomDuration := time.Duration(float64(d) * randomFactor)
	time.Sleep(randomDuration)
}
