package utils

import "time"

const (
	// DateTimeLayout is the layout string for datetime format.
	DateTimeLayout = "2006-01-02 15:04:05"

	// DateTimeLayoutWithMS is the layout string for datetime format with milliseconds.
	DateTimeLayoutWithMS = "2006-01-02 15:04:05.000"

	// RFC3339 is the layout string for RFC3339 format.
	RFC3339 = "2006-01-02T15:04:05Z07:00"

	// DateTimeLayoutWithMSAndTZ is the layout string for datetime format with milliseconds and timezone.
	DateTimeLayoutWithMSAndTZ = "2006-01-02T15:04:05.000Z"

	// TimeLayout is the layout string for time format.
	TimeLayout = "15:04:05"

	// DateLayout is the layout string for date format.
	DateLayout = "2006-01-02"
)

// FormatDateTimeLayout formats the given time to the layout string "2006-01-02 15:04:05".
func FormatDateTimeLayout(t time.Time) string {
	return t.Format(DateTimeLayout)
}

// ParseDateTimeLayout parses the given string to time with layout string "2006-01-02 15:04:05".
func ParseDateTimeLayout(s string) (time.Time, error) {
	return time.Parse(DateTimeLayout, s)
}

// FormatDateTimeLayoutWithMS formats the given time to the layout string "2006-01-02 15:04:05.000".
func FormatDateTimeLayoutWithMS(t time.Time) string {
	return t.Format(DateTimeLayoutWithMS)
}

// ParseDateTimeLayoutWithMS parses the given string to time with layout string "2006-01-02 15:04:05.000".
func ParseDateTimeLayoutWithMS(s string) (time.Time, error) {
	return time.Parse(DateTimeLayoutWithMS, s)
}

// FormatDateTimeRFC3339 formats the given time to the layout string "2006-01-02T15:04:05Z07:00".
func FormatDateTimeRFC3339(t time.Time) string {
	return t.Format(RFC3339)
}

// ParseDateTimeRFC3339 parses the given string to time with layout string "2006-01-02T15:04:05Z07:00".
func ParseDateTimeRFC3339(s string) (time.Time, error) {
	return time.Parse(RFC3339, s)
}

// FormatDateTimeLayoutWithMSAndTZ formats the given time to the layout string "2006-01-02T15:04:05.000Z".
func FormatDateTimeLayoutWithMSAndTZ(t time.Time) string {
	return t.Format(DateTimeLayoutWithMSAndTZ)
}

// ParseDateTimeLayoutWithMSAndTZ parses the given string to time with layout string "2006-01-02T15:04:05.000Z".
func ParseDateTimeLayoutWithMSAndTZ(s string) (time.Time, error) {
	return time.Parse(DateTimeLayoutWithMSAndTZ, s)
}
