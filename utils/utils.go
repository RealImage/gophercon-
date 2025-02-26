package utils

import (
	"time"

	"golang.org/x/time/rate"
)

// Limiter for API requests (slow it down to prevent blocking)
var Limiter = rate.NewLimiter(rate.Every(2*time.Second), 1) // 1 request every 2 seconds
