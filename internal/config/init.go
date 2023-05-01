package config

import (
	"math/rand"
	"time"
)

func init() {
	time.Local = time.UTC
	// for generating pseudorandom numbers
	rand.Seed(time.Now().UnixNano())
}
