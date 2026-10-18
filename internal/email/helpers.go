package email

import "time"

// timeNow returns current UTC time (helper for consistency)
func timeNow() time.Time {
	return time.Now().UTC()
}