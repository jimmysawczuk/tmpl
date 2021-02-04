package tmplfunc

import (
	"bytes"
	"html/template"
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

// SVG reads the file at the provided path and returns its contents under the assumption that it's an HTML-safe string.
func SVG(path string) (template.HTML, error) {
	contents, err := File(path)
	if err != nil {
		return "", err
	}

	return template.HTML(contents), nil
}
