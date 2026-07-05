package pkg

import (
	"fmt"
	"slices"
	"strings"
	"time"

	iso8601duration "github.com/sosodev/duration"
)

// parseDuration parses a duration string like "1y2m" into a time.Duration
func GetTimeDurationFromRelativeDate(s string) (time.Duration, error) {
	// append Prefix "P" to the string to make it ISO 8601 compliant
	if s == "" {
		return 0, fmt.Errorf("duration cannot be empty")
	}

	s = strings.ToUpper("P" + s)

	d, err := iso8601duration.Parse(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return d.ToTimeDuration(), nil
}

func GetMaxDurationFromRelativeDates(dates []string) (time.Duration, error) {
	var durations []time.Duration

	for _, d := range dates {
		dt, err := GetTimeDurationFromRelativeDate(d)
		if err != nil {
			return 0, err
		}
		durations = append(durations, dt)
	}

	return slices.Max(durations), nil
}

func IsValidRelativeDate(s string) bool {
	_, err := GetTimeDurationFromRelativeDate(s)

	if err != nil {
		return false
	}
	return true
}
