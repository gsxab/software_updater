package util

import "time"

func FormatTime(t *time.Time, format string) *string {
	if t == nil {
		return nil
	}
	result := (*t).Format(format)
	return &result
}
