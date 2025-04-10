package utils

import (
	"strings"
	"time"
)

// StartTime stores the timestamp when the server started
var StartTime time.Time

// InitStartTime initializes the StartTime variable to record when the server began running.
func InitStartTime() {
	StartTime = time.Now()
}

// Contains checks if the substring `substr` is present in the string `s`.
// Returns true if the substring exists, false otherwise.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// BoolPtr returns a pointer to the given boolean value `b`.
// Useful for creating optional boolean fields in structs (e.g., for JSON APIs).
func BoolPtr(b bool) *bool {
	return &b
}
