package tmpl

import (
	"bytes"
	"io"
	"os"
	"runtime"
	text "text/template"
	"time"

	"github.com/pkg/errors"
)

type Format int

const (
	FormatText = Format(0)
	FormatJSON = Format(1)
	FormatHTML = Format(2)
)

type goEnv struct {
	OS   string
	Arch string
	Ver  string
}

type Executor interface {
	Execute(io.Writer, interface{}) error
}

type Parser interface {
	Parse(string) (Executor, error)
}

type Tmpl struct {
	Hostname string
	GoEnv    goEnv

	now     time.Time
	envVars map[string]string

	watched []string
}

func New() *Tmpl {
	h, _ := os.Hostname()

	t := &Tmpl{
		Hostname: h,
		GoEnv: goEnv{
			Ver:  runtime.Version(),
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		now: time.Now(),
	}

	return t
}

func (t Tmpl) HTML() *HTMLTmpl {
	return &HTMLTmpl{
		Tmpl: t,
	}
}

func (t Tmpl) Text() *Tmpl {
	return &t
}

func (t Tmpl) JSON() *JSONTmpl {
	return &JSONTmpl{
		Tmpl: t,
	}
}

func (t *Tmpl) WithEnv(m map[string]string) *Tmpl {
	t.envVars = m
	return t
}

func (t *Tmpl) Watch(path string) error {
	t.watched = append(t.watched, path)
	return nil
}

func (t *Tmpl) Execute(out io.Writer, in io.Reader) error {
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
