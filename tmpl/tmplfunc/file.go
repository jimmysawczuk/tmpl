package tmplfunc

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

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

func Link(d FilesystemDepender) func(...string) (template.HTML, error) {
	log.Println("in", d.In().Name())
	log.Println("out", d.Out().Name())
	log.Println("baseDir", d.BaseDir())

	return func(params ...string) (template.HTML, error) {
		if len(params) == 0 {
			return "", errors.Errorf("link requires at least one parameter")
		}

		var localPath string
		var rel string
		var mime string
		var relativePath string

		if len(params) > 0 {
			localPath = params[0]
		}

		if len(params) > 1 {
			rel = params[1]
		}

		if len(params) > 2 {
			mime = params[2]
		}

		if len(params) > 3 {
			relativePath = params[3]
		}

		path, _ := filepath.Rel(d.BaseDir(), localPath)
		if relativePath != "" {
			path = relativePath
		}

		log.Println(d.BaseDir(), localPath, path)

		relAttr := ""
		if rel != "" {
			relAttr = fmt.Sprintf(` rel=%q`, rel)
		}

		mimeAttr := ""
		if mime != "" {
			mimeAttr = fmt.Sprintf(` type=%q`, mime)
		}

		fp, err := os.Open(localPath)
		if err != nil {
			return "", errors.Wrap(err, "os: open")
		}

		defer fp.Close()

		d.Depend(localPath)

		return template.HTML(fmt.Sprintf(`<link href=%q%s%s>`, "/"+path, relAttr, mimeAttr)), nil
	}
}
