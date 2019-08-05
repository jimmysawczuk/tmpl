package main

import (
	"bytes"
	"encoding/json"
	html "html/template"
	"io"
	"io/ioutil"
	text "text/template"

	"github.com/pkg/errors"
)

func read(in io.Reader) string {
	by, _ := ioutil.ReadAll(in)
	return string(by)
}

func writeHTML(in io.Reader, o payload, out io.Writer) error {
	tmpl, err := html.New("output").Funcs(o.tmplfuncs()).Parse(read(in))
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	if err := tmpl.Execute(out, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}

func writeJSON(in io.Reader, o payload, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(read(in))
	if err != nil {
		errors.Wrap(err, "compile template")
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, o); err != nil {
		errors.Wrap(err, "execute template")
	}

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, buf.Bytes()); err != nil {
		return errors.Wrap(err, "compact json")
	}

	io.Copy(out, dst)

	return nil
}

func writeText(in io.Reader, o payload, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(read(in))
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	if err := tmpl.Execute(out, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}
