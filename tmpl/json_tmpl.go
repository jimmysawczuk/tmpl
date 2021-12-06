package tmpl

import (
	"bytes"
	"encoding/json"
	"io"
	text "text/template"

	"github.com/pkg/errors"
)

type JSONTmpl struct {
	*Tmpl

	Minify bool
}

func (t *JSONTmpl) WithMinify(m bool) *JSONTmpl {
	t.Minify = m
	return t
}

func (t *JSONTmpl) Execute(out io.Writer, in io.Reader) error {
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, in); err != nil {
		return errors.Wrap(err, "io: copy (input)")
	}

	tmpl, err := text.New("output").Funcs(t.funcs()).Parse(buf.String())
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf.Reset()

	if err := tmpl.Execute(&buf, t); err != nil {
		errors.Wrap(err, "execute template")
	}

	dst := bytes.Buffer{}
	if t.Minify {
		if err := json.Compact(&dst, buf.Bytes()); err != nil {
			return errors.Wrap(err, "json: compact")
		}
	} else {
		if err := json.Indent(&dst, buf.Bytes(), "", "    "); err != nil {
			return errors.Wrap(err, "json: indent")
		}
	}

	if _, err := io.Copy(out, &dst); err != nil {
		return errors.Wrap(err, "io: copy (output)")
	}

	return nil
}
