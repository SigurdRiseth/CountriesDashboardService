package utils

import (
	"strings"
	"time"
)

// StartTime stores the timestamp when the server started
var StartTime time.Time

// InitStartTime initializes the StartTime variable
func InitStartTime() {
	StartTime = time.Now()
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
