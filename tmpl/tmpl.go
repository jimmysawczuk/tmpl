package tmpl

import (
	"bytes"
	"encoding/json"
	html "html/template"
	"io"
	"os"
	"runtime"
	text "text/template"
	"time"

	"github.com/jimmysawczuk/tmpl/tmpl/tmplfunc"
	"github.com/pkg/errors"
	"github.com/tdewolff/minify"
	htmlminify "github.com/tdewolff/minify/html"
)

type goEnv struct {
	OS   string
	Arch string
	Ver  string
}

type Tmpl struct {
	Hostname string
	GoEnv    goEnv

	now time.Time
}

func New() Tmpl {
	h, _ := os.Hostname()

	return Tmpl{
		Hostname: h,
		GoEnv: goEnv{
			Ver:  runtime.Version(),
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		now: time.Now(),
	}
}

func (t Tmpl) funcs() map[string]interface{} {
	return map[string]interface{}{
		"asset": tmplfunc.AssetLoaderFunc(t.now),
		"env":   tmplfunc.Env,

		"getJSON": tmplfunc.GetJSON,
		"jsonify": tmplfunc.JSONify,

		"now":        tmplfunc.NowFunc(t.now),
		"parseTime":  tmplfunc.ParseTime,
		"formatTime": tmplfunc.FormatTime,
		"timeIn":     tmplfunc.TimeIn,

		"safeHTML":     tmplfunc.SafeHTML,
		"safeHTMLAttr": tmplfunc.SafeAttr,
		"safeJS":       tmplfunc.SafeJS,
		"safeCSS":      tmplfunc.SafeCSS,

		"seq": tmplfunc.Seq,
		"add": tmplfunc.Add,
	}
}

func (t Tmpl) WriteHTML(in io.Reader, out io.Writer, min bool) error {
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, in); err != nil {
		return errors.Wrap(err, "io: copy (input)")
	}

	tmpl, err := html.New("output").Funcs(t.funcs()).Parse(buf.String())
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf.Reset()

	if err := tmpl.Execute(&buf, t); err != nil {
		return errors.Wrap(err, "execute template")
	}

	if min {
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

func (t Tmpl) WriteJSON(in io.Reader, out io.Writer, min bool) error {
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
	if min {
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

func (t Tmpl) WriteText(in io.Reader, out io.Writer) error {
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, in); err != nil {
		return errors.Wrap(err, "io: copy (input)")
	}

	tmpl, err := text.New("output").Funcs(t.funcs()).Parse(buf.String())
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf.Reset()

	if err := tmpl.Execute(out, t); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}
