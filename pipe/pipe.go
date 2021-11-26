package pipe

import (
	"io"
	"os"

	"github.com/jimmysawczuk/tmpl/tmpl"
	"github.com/pkg/errors"
)

type Pipe struct {
	In      string
	Out     string
	BaseDir string

	Format string
	Mode   tmpl.Mode

	Minify bool
	Env    map[string]string
	Params map[string]interface{}

	deps []string
}

type executor interface {
	Execute(io.Writer, io.Reader) error
	Dependencies() []string
}

func (p *Pipe) Run() error {
	in, err := os.Open(p.In)
	if err != nil {
		return errors.Wrapf(err, "open input (path: %s)", p.In)
	}
	defer in.Close()

	out, err := os.OpenFile(p.Out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open output (path: %s)", p.Out)
	}
	defer out.Close()

	var t executor
	switch p.Format {
	case "html":
		t = tmpl.New().WithMode(p.Mode).WithBaseDir(p.BaseDir).WithIO(in, out).HTML().WithMinify(p.Minify)
	case "json":
		t = tmpl.New().WithMode(p.Mode).WithBaseDir(p.BaseDir).WithIO(in, out).JSON().WithMinify(p.Minify)
	default:
		t = tmpl.New().WithMode(p.Mode).WithBaseDir(p.BaseDir).WithIO(in, out)
	}

	if err := t.Execute(out, in); err != nil {
		return errors.Wrapf(err, "execute (%T, in: %s)", t, in)
	}

	// TODO: pipe dependencies up to refs

	return nil
}
