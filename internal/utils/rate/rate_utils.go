// Package rate provides rate limiting utilities
// Author: Done-0
// Created: 2025-08-31
package rate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// ParseLimit parses rate limit string like "60/min", "1/s"
func ParseLimit(s string) (rate.Limit, int, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: %s", s)
	}

	requests, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid requests: %s", parts[0])
	}

	var duration time.Duration
	switch parts[1] {
	case "s", "sec", "second":
		duration = time.Second
	case "m", "min", "minute":
		duration = time.Minute
	case "h", "hour":
		duration = time.Hour
	default:
		duration, err = time.ParseDuration(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid duration: %s", parts[1])
		}
	}

	return rate.Every(duration / time.Duration(requests)), requests, nil
}
