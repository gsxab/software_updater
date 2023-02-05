package util

import "time"

func FormatTime(t *time.Time, format string) *string {
	if t == nil {
		return nil
	}
	result := (*t).Format(format)
	return &result
}

func FormatTimeInt64(t int64, format string) string {
	result := time.Unix(t, 0).Format(format)
	return result
}
