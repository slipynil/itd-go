package root

import (
	"time"
)

type Config struct {
	RefreshToken string
	Url          string
	Domain       string
	UserAgent    string
	Timeout      time.Duration
}
