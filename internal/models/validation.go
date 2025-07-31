package models

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidTimezone   = errors.New("invalid timezone")
	ErrInvalidColor      = errors.New("invalid color format")
	ErrInvalidTimeRange  = errors.New("invalid time range")
	ErrInvalidDateFormat = errors.New("invalid date format")
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

// ValidateActivityType validates activity type parameter
func ValidateActivityType(activityType string) bool {
	if activityType == "" {
		return true // Optional parameter
	}
	return ActivityType(activityType).IsValid()
}

// ValidateValueRating validates value rating parameter
func ValidateValueRating(valueRating string) bool {
	if valueRating == "" {
		return true // Optional parameter
	}
	return ValueRating(valueRating).IsValid()
}

// ValidateImpactLevel validates impact level parameter
func ValidateImpactLevel(impactLevel string) bool {
	if impactLevel == "" {
		return true // Optional parameter
	}
	return ImpactLevel(impactLevel).IsValid()
}

// ValidateDateFormat validates date string format (YYYY-MM-DD)
func ValidateDateFormat(dateStr string) error {
	if dateStr == "" {
		return nil // Optional parameter
	}
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ErrInvalidDateFormat
	}
	return nil
}
