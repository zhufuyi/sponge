package utils

import "time"

const (
	// DateTimeLayout is the layout string for datetime format.
	DateTimeLayout = "2006-01-02 15:04:05"

	// DateTimeLayoutWithMS is the layout string for datetime format with milliseconds.
	DateTimeLayoutWithMS = "2006-01-02 15:04:05.000"

	// DateTimeLayoutWithMSAndTZ is the layout string for datetime format with milliseconds and timezone.
	DateTimeLayoutWithMSAndTZ = "2006-01-02T15:04:05.000Z"

	// TimeLayout is the layout string for time format.
	TimeLayout = "15:04:05"

	// DateLayout is the layout string for date format.
	DateLayout = "2006-01-02"
)

// FormatDateTime formats the given time to string
func FormatDateTime(t time.Time, format string) string {
	switch format {
	case DateTimeLayoutWithMS:
		return t.Format(DateTimeLayoutWithMS)
	case DateTimeLayoutWithMSAndTZ:
		return t.UTC().Format(DateTimeLayoutWithMSAndTZ)
	case TimeLayout:
		return t.Format(TimeLayout)
	case DateLayout:
		return t.Format(DateLayout)
	default:
		return t.Format(DateTimeLayout)
	}
}

// ParseDateTime parses the given string to time
func ParseDateTime(s string, format string) (time.Time, error) {
	switch format {
	case DateTimeLayoutWithMS:
		return time.Parse(DateTimeLayoutWithMS, s)
	case DateTimeLayoutWithMSAndTZ:
		return time.Parse(DateTimeLayoutWithMSAndTZ, s)
	case TimeLayout:
		return time.Parse(TimeLayout, s)
	case DateLayout:
		return time.Parse(DateLayout, s)
	default:
		return time.Parse(DateTimeLayout, s)
	}
}
