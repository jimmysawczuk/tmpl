package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/jimmysawczuk/tmpl/tmplfunc"
	"github.com/pkg/errors"
)

type goEnv struct {
	OS   string
	Arch string
	Ver  string
}

type payload struct {
	Hostname string
	GoEnv    goEnv

	now  time.Time
	mode string
}

func getHostname() string {
	v, _ := os.Hostname()
	return v
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(2)
}

func main() {
	o := payload{
		Hostname: getHostname(),
		GoEnv: goEnv{
			Ver:  runtime.Version(),
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},

		mode: os.Getenv("MODE"),
		now:  time.Now(),
	}

	tmplStr, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fatalErr(errors.Wrapf(err, "read template file: %s", os.Args[1]))
	}

	tmpl, err := template.New("output").Funcs(map[string]interface{}{
		"asset": tmplfunc.AssetLoaderFunc(o.now, o.mode),
		"env":   tmplfunc.Env,

		"getJSON": tmplfunc.GetJSON,
		"jsonify": tmplfunc.JSONify,

		"now":        tmplfunc.NowFunc(o.now),
		"parseTime":  tmplfunc.ParseTime,
		"formatTime": tmplfunc.FormatTime,

		"safeHTML":     tmplfunc.SafeHTML,
		"safeHTMLAttr": tmplfunc.SafeAttr,
		"safeJS":       tmplfunc.SafeJS,
		"safeCSS":      tmplfunc.SafeCSS,
	}).Parse(string(tmplStr))
	if err != nil {
		fatalErr(errors.Wrapf(err, "compile template: %s", os.Args[1]))
	}

	err = tmpl.Execute(os.Stdout, o)
	if err != nil {
		fatalErr(errors.Wrapf(err, "execute template: %s", err))
	}
}
