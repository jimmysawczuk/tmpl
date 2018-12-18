package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var config = struct {
	Output          string
	Format          string
	EnvFile         string
	TimestampAssets bool
}{
	TimestampAssets: true,
}

func init() {
	flag.StringVar(&config.Output, "o", "", "output destination ('' means stdout)")
	flag.StringVar(&config.Format, "fmt", "", "output format ('', 'html' or 'json')")
	flag.StringVar(&config.EnvFile, "env-file", "", "pipe .env file before executing template")
	flag.BoolVar(&config.TimestampAssets, "timestamp-assets", true, "set to false to not automatically timestamp production assets")
}

func main() {
	o, err := newPayload()
	if err != nil {
		fatalErr(errors.Wrap(err, "build payload"))
	}

	flag.Parse()

	if config.EnvFile != "" {
		if err := godotenv.Load(config.EnvFile); err != nil {
			fatalErr(errors.Wrapf(err, "load .env file %s", config.EnvFile))
		}
	}

	by, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fatalErr(errors.Wrapf(err, "read template file: %s", flag.Arg(0)))
	}

	var out io.Writer = os.Stdout
	if config.Output != "" {
		fp, err := openOutputFile(config.Output)
		if err != nil {
			fatalErr(errors.Wrapf(err, "open output file %s"))
		} else {
			out = fp
		}

		out = fp
	}

	switch config.Format {
	case "html":
		if err := writeHTML(string(by), o, out); err != nil {
			fatalErr(err)
		}
	case "json":
		if err := writeJSON(string(by), o, out); err != nil {
			fatalErr(err)
		}
	default:
		if err := writeText(string(by), o, out); err != nil {
			fatalErr(err)
		}
	}
}

func openOutputFile(path string) (io.Writer, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, errors.Wrap(err, "mkdir")
	}

	fp, err := os.OpenFile(config.Output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok {
			return nil, errors.New("unexpected error; couldn't assert to *os.PathError")
		}

		return nil, perr
	}

	return fp, nil
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(2)
}
