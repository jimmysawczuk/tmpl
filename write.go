package main

import (
	"bytes"
	"encoding/json"
	html "html/template"
	"io"
	text "text/template"

	"github.com/pkg/errors"
)

func writeHTML(in string, o payload, out io.Writer) error {
	tmpl, err := html.New("output").Funcs(o.tmplfuncs()).Parse(in)
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	if err := tmpl.Execute(out, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}

func writeJSON(in string, o payload, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(in)
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

func writeText(in string, o payload, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(in)
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	if err := tmpl.Execute(out, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}
