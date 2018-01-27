package tmplfunc

import (
	"time"
)

func NowFunc(now time.Time) func() time.Time {
	return func() time.Time { return now }
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func FormatTime(format string, t time.Time) string {
	return t.Format(format)
}
