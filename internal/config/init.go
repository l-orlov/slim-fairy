package config

import "time"

func init() {
	time.Local = time.UTC
}
