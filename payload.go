package main

import (
	"os"
	"runtime"
	"time"

	"github.com/jimmysawczuk/tmpl/tmplfunc"
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

func newPayload() (payload, error) {
	h, _ := os.Hostname()

	return payload{
		Hostname: h,
		GoEnv: goEnv{
			Ver:  runtime.Version(),
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},

		mode: os.Getenv("MODE"),
		now:  time.Now(),
	}, nil
}

func (o payload) tmplfuncs() map[string]interface{} {
	return map[string]interface{}{
		"asset": tmplfunc.AssetLoaderFunc(o.now, o.mode == "production" && config.TimestampAssets),
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

		"seq": tmplfunc.Seq,
	}
}
