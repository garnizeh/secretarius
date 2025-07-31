package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper functions for pgtype conversions

// stringRepeat creates a string by repeating the input string a specified number of times
func stringRepeat(s string, count int) string {
	result := ""
	for range count {
		result += s
	}
	return result
}

// stringPtr creates a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// uuidToPgUUID converts *uuid.UUID to pgtype.UUID
func uuidToPgUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

// stringToPgText converts *string to pgtype.Text
func stringToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// timeToPgTimestamptz converts time.Time to pgtype.Timestamptz
func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// pgUUIDToUUID converts pgtype.UUID to *uuid.UUID
func pgUUIDToUUID(pgUUID pgtype.UUID) *uuid.UUID {
	if !pgUUID.Valid {
		return nil
	}
	u := uuid.UUID(pgUUID.Bytes)
	return &u
}

// pgTextToString converts pgtype.Text to *string
func pgTextToString(pgText pgtype.Text) *string {
	if !pgText.Valid {
		return nil
	}
	return &pgText.String
}

// pgTimestamptzToTime converts pgtype.Timestamptz to time.Time
func pgTimestamptzToTime(pgTime pgtype.Timestamptz) time.Time {
	if !pgTime.Valid {
		return time.Time{}
	}
	return pgTime.Time
}

// pgInt4ToInt converts pgtype.Int4 to int
func pgInt4ToInt(pgInt pgtype.Int4) int {
	if !pgInt.Valid {
		return 0
	}
	return int(pgInt.Int32)
}

// pgTextToStringRequired converts pgtype.Text to string for required fields
func pgTextToStringRequired(pgText pgtype.Text) string {
	if !pgText.Valid {
		return ""
	}
	return pgText.String
}

// stringToPgTextRequired converts string to pgtype.Text for required fields
func stringToPgTextRequired(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// timeToPgDate converts *time.Time to pgtype.Date
func timeToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

// pgDateToTime converts pgtype.Date to *time.Time
func pgDateToTime(pgDate pgtype.Date) *time.Time {
	if !pgDate.Valid {
		return nil
	}
	return &pgDate.Time
}

// boolToPgBool converts bool to pgtype.Bool
func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

// pgBoolToBool converts pgtype.Bool to bool
func pgBoolToBool(pgBool pgtype.Bool) bool {
	if !pgBool.Valid {
		return false
	}
	return pgBool.Bool
}

// Helper function to convert pgtype.UUID to *uuid.UUID
func pgUUIDToUUIDPtr(pgUUID pgtype.UUID) *uuid.UUID {
	if !pgUUID.Valid {
		return nil
	}
	id := uuid.UUID(pgUUID.Bytes)
	return &id
}

// Helper function to convert pgtype.Int4 to int32
func pgInt4ToInt32(pgInt pgtype.Int4) int32 {
	if !pgInt.Valid {
		return 0
	}
	return pgInt.Int32
}
