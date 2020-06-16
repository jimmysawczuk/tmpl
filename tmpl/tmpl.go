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

	now     time.Time
	envVars map[string]string
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

func (t *Tmpl) WithEnv(m map[string]string) *Tmpl {
	t.envVars = m
	return t
}

func (t *Tmpl) funcs() map[string]interface{} {
	return map[string]interface{}{
		"add":          tmplfunc.Add,
		"asset":        tmplfunc.Asset,
		"env":          tmplfunc.EnvFunc(t.envVars),
		"formatTime":   tmplfunc.FormatTime,
		"getJSON":      tmplfunc.GetJSON,
		"jsonify":      tmplfunc.JSONify,
		"now":          tmplfunc.NowFunc(t.now),
		"parseTime":    tmplfunc.ParseTime,
		"safeCSS":      tmplfunc.SafeCSS,
		"safeHTML":     tmplfunc.SafeHTML,
		"safeHTMLAttr": tmplfunc.SafeAttr,
		"safeJS":       tmplfunc.SafeJS,
		"seq":          tmplfunc.Seq,
		"sub":          tmplfunc.Sub,
		"timeIn":       tmplfunc.TimeIn,
	}
}

func (t Tmpl) WriteHTML(out io.Writer, in io.Reader, min bool) error {
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

func (t Tmpl) WriteJSON(out io.Writer, in io.Reader, min bool) error {
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

func (t Tmpl) WriteText(out io.Writer, in io.Reader) error {
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
