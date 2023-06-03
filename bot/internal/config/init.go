package config

import (
	"math/rand"
	"time"
)

// RandomGenerator is used for generating pseudorandom numbers
var RandomGenerator *rand.Rand

func init() {
	time.Local = time.UTC
	RandomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
}
