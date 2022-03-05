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

type goEnv struct {
	OS   string
	Arch string
	Ver  string
}

type Mode int

const (
	ModeLocal      Mode = 0
	ModeProduction Mode = iota
)

type Tmpl struct {
	Hostname string
	GoEnv    goEnv

	mode    Mode
	in      *os.File
	out     *os.File
	baseDir string

	leftDelim  string
	rightDelim string

	now     time.Time
	envVars map[string]string

	refs map[string]struct{}
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
		leftDelim:  "{{",
		rightDelim: "}}",
		now:        time.Now(),
		refs:       map[string]struct{}{},
	}

	return t
}

func (t *Tmpl) WithIO(in, out *os.File) *Tmpl {
	t.in = in
	t.out = out
	return t
}

func (t *Tmpl) WithBaseDir(dir string) *Tmpl {
	t.baseDir = dir
	return t
}

func (t *Tmpl) WithMode(mode Mode) *Tmpl {
	t.mode = mode
	return t
}

func (t *Tmpl) WithDelims(left, right string) *Tmpl {
	if left != "" && right != "" {
		t.leftDelim = left
		t.rightDelim = right
	}
	return t
}

func (t *Tmpl) In() *os.File {
	return t.in
}

func (t *Tmpl) Out() *os.File {
	return t.out
}

func (t *Tmpl) BaseDir() string {
	return t.baseDir
}

func (t *Tmpl) IsProduction() bool {
	return t.mode == ModeProduction
}

func (t *Tmpl) HTML() *HTMLTmpl {
	return &HTMLTmpl{
		Tmpl: t,
	}
}

func (t *Tmpl) JSON() *JSONTmpl {
	return &JSONTmpl{
		Tmpl: t,
	}
}

func (t *Tmpl) WithEnv(m map[string]string) *Tmpl {
	t.envVars = m
	return t
}

func (t *Tmpl) Ref(path string) error {
	t.refs[path] = struct{}{}
	return nil
}

func (t *Tmpl) Refs() []string {
	tbr := make([]string, len(t.refs))
	i := 0
	for k := range t.refs {
		tbr[i] = k
		i++
	}
	return tbr
}

func (t *Tmpl) Execute(out io.Writer, in io.Reader) error {
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, in); err != nil {
		return errors.Wrap(err, "io: copy (input)")
	}

	tmpl, err := text.New("output").Funcs(t.funcs()).Delims(t.leftDelim, t.rightDelim).Parse(buf.String())
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf.Reset()

	if err := tmpl.Execute(out, t); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}
