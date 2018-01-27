package tmplfunc

import (
	"fmt"
	"regexp"
	"time"
)

var timestampRegexp = regexp.MustCompile(`\.(\w+)$`)

func AssetLoaderFunc(now time.Time, mode string) func(path string) string {
	return func(path string) string {
		if mode == "production" {
			return timestampRegexp.ReplaceAllString(path, fmt.Sprintf(".%d.$1", now.Unix()))
		}

		return path
	}
}
