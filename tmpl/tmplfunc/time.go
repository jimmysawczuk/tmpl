package tmplfunc

import (
	"time"
)

// NowFunc returns a function which returns the provided time.Time.
func NowFunc(now time.Time) func() time.Time {
	return func() time.Time { return now }
}

// ParseTime parses the provided string using the time.RFC3339 format. It
// will return an error if Parse returns an error.
//
// See: https://golang.org/pkg/time#Parse, https://golang.org/pkg/time#RFC3339
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// FormatTime returns the provided time formatted using the provided format string.
func FormatTime(format string, t time.Time) string {
	return t.Format(format)
}

// TimeIn returns the provided time with the location set to the provided location
// string. TimeIn will return an error if tz doesn't correspond to a valid
// location string.
//
// See: https://golang.org/pkg/time/#LoadLocation
func TimeIn(tz string, t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return t, err
	}

	return t.In(loc), nil
}
