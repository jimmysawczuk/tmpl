package tmplfunc

import (
	"os"
	"strings"
)

// EnvFunc returns a function which first checks the provided map for the requested
// key, falling back to os.Getenv if the key is not present in the map.
func EnvFunc(m map[string]string) func(string) string {
	return func(s string) string {
		if v, ok := m[strings.ToLower(s)]; ok {
			return v
		}

		return os.Getenv(s)
	}
}
