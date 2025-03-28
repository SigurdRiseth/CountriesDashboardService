package utils

import "time"

// StartTime stores the timestamp when the server started
var StartTime time.Time

// InitStartTime initializes the StartTime variable
func InitStartTime() {
	StartTime = time.Now()
}
