package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func IsValidRelativeDate(s string) bool {
	if s == "" {
		return false
	}
	_, err := ParseDuration(s)

	if err != nil {
		return false
	}
	return true
}

func IsValidDay(day string) bool {
	switch strings.ToLower(day) {
	case "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday":
		return true
	default:
		return false
	}
}

// Regex to validate and capture tokens like "2w", "10d", "5h", "30m"
var customDurationRegex = regexp.MustCompile(`^(\d+)([wdhm])$`)

// ParseCustomDuration converts AFK duration string format into a time.Duration
func ParseDuration(input string) (time.Duration, error) {
	// Clean up leading/trailing whitespaces
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, errors.New("duration string cannot be empty")
	}

	// Split by any amount of whitespace (handles single space, multiple spaces, tabs)
	tokens := strings.Fields(input)
	var totalDuration time.Duration

	for _, token := range tokens {
		matches := customDurationRegex.FindStringSubmatch(token)
		if len(matches) != 3 {
			return 0, fmt.Errorf("invalid duration component: %q (must be a number followed by [w]eeks, [d]ays, [h]ours, or [m]inutes)", token)
		}

		// Extract the numeric value (e.g., "30" from "30m")
		value, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number value in component %q: %w", token, err)
		}

		// Map the single-letter unit to Go's time.Duration
		var unitDuration time.Duration
		switch matches[2] {
		case "w": // Weeks
			unitDuration = time.Hour * 24 * 7
		case "d": // Days
			unitDuration = time.Hour * 24
		case "h": // Hours
			unitDuration = time.Hour
		case "m": // Mins
			unitDuration = time.Minute
		}

		// Accumulate the duration values safely
		totalDuration += time.Duration(value) * unitDuration
	}

	return totalDuration, nil
}
