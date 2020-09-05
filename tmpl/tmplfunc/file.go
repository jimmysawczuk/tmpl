package tmplfunc

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
)

// File reads the file at the provided path and returns its contents as a string.
func File(path string) (string, error) {
	fp, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "os: open")
	}

	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, fp); err != nil {
		return "", errors.Wrap(err, "io: copy")
	}

	return buf.String(), nil
}
