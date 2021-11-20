package tmpl

import (
	"bytes"
	html "html/template"
	"io"

	"github.com/pkg/errors"
	"github.com/tdewolff/minify"
	htmlminify "github.com/tdewolff/minify/html"
)

type HTMLTmpl struct {
	*Tmpl

	Minify bool
}

func (t *HTMLTmpl) WithMinify(m bool) *HTMLTmpl {
	t.Minify = m
	return t
}

func (t *HTMLTmpl) Execute(out io.Writer, in io.Reader) error {
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, in); err != nil {
		return errors.Wrap(err, "io: copy (input)")
	}

	tmpl, err := html.New("output").Funcs(t.Tmpl.funcs()).Parse(buf.String())
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf.Reset()

	if err := tmpl.Execute(&buf, t); err != nil {
		return errors.Wrap(err, "execute template")
	}

	if t.Minify {
		by := buf.Bytes()
		buf.Reset()

		m := minify.New()
		hm := htmlminify.DefaultMinifier
		hm.KeepDocumentTags = true

		hm.Minify(m, &buf, bytes.NewReader(by), nil)
	}

	if _, err := io.Copy(out, &buf); err != nil {
		return errors.Wrap(err, "io: copy (output)")
	}

	return nil
}
