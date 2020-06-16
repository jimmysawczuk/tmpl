package tmplfunc

import (
	"os"
	"strings"
)

func EnvFunc(m map[string]string) func(string) string {
	return func(s string) string {
		if v, ok := m[strings.ToLower(s)]; ok {
			return v
		}

		return os.Getenv(s)
	}
}
