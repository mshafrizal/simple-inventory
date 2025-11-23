package util

import "time"

func FormatTimeUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
