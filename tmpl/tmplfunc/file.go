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

// Inline reads the file at the provided path and returns its contents as a string. It
// also marks the file as updateable.
func Inline(w Watcher) func(string) (string, error) {
	return func(path string) (string, error) {
		fp, err := os.Open(path)
		if err != nil {
			return "", errors.Wrap(err, "os: open")
		}

		w.Watch(path)

		buf := bytes.Buffer{}
		if _, err := io.Copy(&buf, fp); err != nil {
			return "", errors.Wrap(err, "io: copy")
		}

		return buf.String(), nil
	}
}

// SVG reads the file at the provided path and returns its contents under the assumption
// that it's an HTML-safe string. It also marks the SVG as updateable.
func SVG(w Watcher) func(string) (template.HTML, error) {
	return func(path string) (template.HTML, error) {
		res, err := Inline(w)(path)
		if err != nil {
			return "", errors.Wrap(err, "inline")
		}

		return template.HTML(res), nil
	}
}
