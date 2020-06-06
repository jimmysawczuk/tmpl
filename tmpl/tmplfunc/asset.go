package tmplfunc

import (
	"time"
)

func AssetLoaderFunc(now time.Time) func(path string) string {
	return func(path string) string {
		return path
	}
}
