package tmplfunc

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

// Asset returns the path provided. In the future, Asset may gain
// the ability to clean or validate the path.
func Asset(path string) string {
	return path
}

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
func Inline(d Depender) func(string) (string, error) {
	return func(path string) (string, error) {
		fp, err := os.Open(path)
		if err != nil {
			return "", errors.Wrap(err, "os: open")
		}

		defer fp.Close()

		d.Depend(path)

		buf := bytes.Buffer{}
		if _, err := io.Copy(&buf, fp); err != nil {
			return "", errors.Wrap(err, "io: copy")
		}

		return buf.String(), nil
	}
}

// SVG reads the file at the provided path and returns its contents under the assumption
// that it's an HTML-safe string. It also marks the SVG as updateable.
func SVG(d Depender) func(string) (template.HTML, error) {
	return func(path string) (template.HTML, error) {
		res, err := Inline(d)(path)
		if err != nil {
			return "", errors.Wrap(err, "inline")
		}

		return template.HTML(res), nil
	}
}

// Ref marks the provided file as a dependency of the template, so any changes to that file
// will trigger a rebuild. It returns no output.
func Ref(d Depender) func(string) string {
	return func(filePath string) string {
		if _, err := os.Stat(filePath); err != nil {
			log.Printf("ref: os: stat: %s", err)
		}

		d.Depend(filePath)
		return ""
	}
}
