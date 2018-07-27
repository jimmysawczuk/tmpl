package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var config struct {
	Output string
	Format string
}

func init() {
	flag.StringVar(&config.Output, "o", "", "output destination ('' means stdout)")
	flag.StringVar(&config.Format, "fmt", "", "output format ('', 'html' or 'json')")
}

func main() {
	o, err := newPayload()
	if err != nil {
		fatalErr(errors.Wrap(err, "build payload"))
	}

	flag.Parse()

	by, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fatalErr(errors.Wrapf(err, "read template file: %s", flag.Arg(0)))
	}

	var out io.Writer = os.Stdout
	if config.Output != "" {
		fp, err := os.Open(config.Output)
		if err != nil {
			fatalErr(errors.Wrapf(err, "open output file %s", config.Output))
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

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(2)
}
