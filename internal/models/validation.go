package models

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidTimezone  = errors.New("invalid timezone")
	ErrInvalidColor     = errors.New("invalid color format")
	ErrInvalidTimeRange = errors.New("invalid time range")
)

// ValidateTimezone validates if the provided timezone string is valid
func ValidateTimezone(tz string) error {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return ErrInvalidTimezone
	}
	return nil
}

// ValidateHexColor validates if the provided color string is a valid hex color
func ValidateHexColor(color string) error {
	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, color)
	if !matched {
		return ErrInvalidColor
	}
	return nil
}

// ValidateTimeRange validates if the end time is after the start time
func ValidateTimeRange(start, end time.Time) error {
	if end.Before(start) {
		return ErrInvalidTimeRange
	}
	return nil
}
