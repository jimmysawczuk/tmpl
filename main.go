package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/jimmysawczuk/tmpl/tmpl"
	"github.com/pkg/errors"
)

type Block struct {
	In     string `json:"in"`
	Out    string `json:"out"`
	Format string `json:"format"`

	Options BlockOpts `json:"options"`
}

type BlockOpts struct {
	Minify bool              `json:"minify"`
	Env    map[string]string `json:"env"`
}

var blocks []Block
var watchMode bool
var configFile string

type pipeline struct {
	inpath  string
	outpath string
	format  string

	minify bool
	env    map[string]string
}

func init() {
	flag.StringVar(&configFile, "f", "./tmpl.config.json", "path to tmpl config file")
	flag.BoolVar(&watchMode, "w", false, "run in watch mode")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}
}

func run() error {
	fp, err := os.Open(configFile)
	if err != nil {
		return errors.Wrapf(err, "open config file (path: %s)", configFile)
	}

	if err := json.NewDecoder(fp).Decode(&blocks); err != nil {
		return errors.Wrap(err, "json: decode config file")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "fsnotify: new")
	}
	defer watcher.Close()

	pipelines := []pipeline{}

	for i, b := range blocks {

		pipe := pipeline{
			format: b.Format,
			minify: b.Options.Minify,
			env:    b.Options.Env,
		}

		pipe.inpath, _ = filepath.Abs(b.In)

		ostat, err := os.Stat(b.Out)
		if err == nil {
			if ostat.IsDir() {
				outpath := filepath.Join(b.Out, filepath.Base(b.In))

				_, err := os.Stat(outpath)
				if err == nil || os.IsNotExist(err) {
					pipe.outpath, _ = filepath.Abs(outpath)
				} else {
					return errors.Wrapf(err, "stat output path (path: %s, block: %d)", outpath, i+1)
				}

				pipe.outpath, _ = filepath.Abs(filepath.Join(b.Out, filepath.Base(b.In)))
			} else {
				pipe.outpath, _ = filepath.Abs(b.Out)
			}
		} else if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(b.Out), 0755); err != nil {
				return errors.Wrapf(err, "mkdir (path: %s, block: %d)", filepath.Dir(b.Out), i+1)
			}

			pipe.outpath, _ = filepath.Abs(b.Out)
		} else {
			return errors.Wrapf(err, "stat (output: %s, block: %d)", b.Out, i+1)
		}

		pipelines = append(pipelines, pipe)

		if watchMode {
			log.Println("watching", pipe.inpath)
			watcher.Add(pipe.inpath)
		}
	}

	for _, pipe := range pipelines {
		if err := pipe.fire(); err != nil {
			return errors.Wrapf(err, "pipeline (path: %s)", pipe.inpath)
		}

	}

	if watchMode {
		done := make(chan bool)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("changed:", event.Name)

						for _, pipe := range pipelines {
							if pipe.inpath == event.Name {
								if err := pipe.fire(); err != nil {
									log.Printf("%s", errors.Wrapf(err, "pipeline (path: %s)", pipe.inpath))
								}
								log.Println(" --> wrote:", pipe.outpath)
							}
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()
		<-done
	}

	return nil
}

func (p *pipeline) fire() error {

	in, err := os.Open(p.inpath)
	if err != nil {
		return errors.Wrapf(err, "open input (path: %s)", p.inpath)
	}
	defer in.Close()

	out, err := os.OpenFile(p.outpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open output (path: %s)", p.outpath)
	}
	defer out.Close()

	tmpl := tmpl.New()

	switch p.format {
	case "html":
		if err := tmpl.WriteHTML(in, out, p.minify); err != nil {
			return errors.Wrapf(err, "write html (in: %s)", p.inpath)
		}
	case "json":
		if err := tmpl.WriteJSON(in, out, p.minify); err != nil {
			return errors.Wrapf(err, "write json (in: %s)", p.inpath)
		}
	default:
		if err := tmpl.WriteText(in, out); err != nil {
			return errors.Wrapf(err, "write text (in: %s)", p.inpath)
		}
	}

	return nil
}
