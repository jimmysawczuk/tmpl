package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jimmysawczuk/tmpl/tmpl"
	"github.com/pkg/errors"
)

var (
	version  string = "dev"
	revision string
	date     string = time.Now().Format(time.RFC3339)
)

var (
	watchMode   bool
	serverMode  bool
	port        int
	baseDir     string
	configFile  string
	showVersion bool
	runCommand  []string
)

type pipeline struct {
	inpath  string
	outpath string
	format  string

	mode   tmpl.Mode
	minify bool
	env    map[string]string

	dependencies []string
}

type block struct {
	In     string `json:"in"`
	Out    string `json:"out"`
	Format string `json:"format"`

	Options blockOpts `json:"options"`
}

type blockOpts struct {
	Minify bool              `json:"minify"`
	Env    map[string]string `json:"env"`
}

var blocks []block

type Tmpl interface {
	Execute(io.Writer, io.Reader) error
	Dependencies() []string
}

func init() {
	flag.Usage = func() {
		fmt.Printf("tmpl %s\n\n", version)

		fmt.Printf("Usage:\n")
		fmt.Printf("  tmpl [options] [-- command]\n\n")

		flag.PrintDefaults()
	}

	flag.StringVar(&configFile, "f", "./tmpl.config.json", "path to tmpl config file")
	flag.BoolVar(&watchMode, "w", false, "run in watch mode")
	flag.BoolVar(&serverMode, "s", false, "run in watch mode and serve")
	flag.IntVar(&port, "p", 8080, "port to listen on in serve mode")
	flag.StringVar(&baseDir, "dir", ".", "public dir")
	flag.BoolVar(&showVersion, "v", false, "show version information")
}

func main() {
	flag.Parse()

	if showVersion {
		flag.Usage()
		os.Exit(0)
	}

	if args := flag.Args(); len(args) > 0 {
		runCommand = args
	}

	watchMode = watchMode || serverMode

	if err := run(); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}
}

func run() error {
	var dependencyMap map[string]map[string]bool = map[string]map[string]bool{}

	fp, err := os.Open(configFile)
	if err != nil {
		return errors.Wrapf(err, "open config file (path: %s)", configFile)
	}

	if err := json.NewDecoder(fp).Decode(&blocks); err != nil {
		return errors.Wrap(err, "json: decode config file")
	}

	var watcher *fsnotify.Watcher
	if watchMode {
		w, err := fsnotify.NewWatcher()
		if err != nil {
			return errors.Wrap(err, "fsnotify: new")
		}

		watcher = w
		defer watcher.Close()
	}

	mode := tmpl.ModeProduction
	if watchMode || serverMode {
		mode = tmpl.ModeLocal
	}

	pipelines := []pipeline{}

	for i, b := range blocks {

		pipe := pipeline{
			format: b.Format,
			minify: b.Options.Minify,
			env:    b.Options.Env,
			mode:   mode,
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
		if err := pipe.run(); err != nil {
			return errors.Wrapf(err, "run pipeline (path: %s)", pipe.inpath)
		}

		if watchMode {
			for _, dep := range pipe.dependencies {
				abs, _ := filepath.Abs(dep)

				watcher.Add(abs)

				if dependencyMap[pipe.inpath] == nil {
					dependencyMap[pipe.inpath] = map[string]bool{}
				}

				dependencyMap[pipe.inpath][abs] = true
			}
		}
	}

	var cmd *exec.Cmd
	if len(runCommand) > 0 {
		cmd = exec.Command(runCommand[0], runCommand[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		if err := cmd.Start(); err != nil {
			return errors.Wrapf(err, "command: start (%s)", strings.Join(runCommand, " "))
		}
	}

	watcherCh := make(chan string)
	if serverMode {
		log.Printf("starting server on :%d", port)

		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.Dir(baseDir)))
		mux.Handle("/__tmpl", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			res := <-watcherCh
			json.NewEncoder(w).Encode(struct {
				Status string    `json:"status"`
				File   string    `json:"file"`
				Date   time.Time `json:"date"`
			}{
				Status: "OK",
				File:   res,
				Date:   time.Now().Truncate(time.Millisecond),
			})
		}))

		go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
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
								if err := pipe.run(); err != nil {
									log.Printf("%s", errors.Wrapf(err, "pipeline (path: %s)", pipe.inpath))
								}
								log.Println(" --> wrote:", pipe.outpath)
							} else if dependencyMap[pipe.inpath] != nil && dependencyMap[pipe.inpath][event.Name] {
								if err := pipe.run(); err != nil {
									log.Printf("%s", errors.Wrapf(err, "pipeline (path: %s)", pipe.inpath))
								}
								log.Println(" --> wrote:", pipe.outpath)
							}
						}
					}

					select {
					case watcherCh <- event.Name:
						log.Println(" ! notified listener")
					default:
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
	} else if cmd != nil {
		err := cmd.Wait()
		if err != nil {
			return errors.Wrapf(err, "command: wait (%s)", strings.Join(runCommand, " "))
		}
	}

	return nil
}

func (p *pipeline) run() error {
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

	var t Tmpl
	switch p.format {
	case "html":
		t = tmpl.New().WithMode(p.mode).WithBaseDir(baseDir).WithIO(in, out).HTML().WithMinify(p.minify)
	case "json":
		t = tmpl.New().WithMode(p.mode).WithBaseDir(baseDir).WithIO(in, out).JSON().WithMinify(p.minify)
	default:
		t = tmpl.New().WithMode(p.mode).WithBaseDir(baseDir).WithIO(in, out)
	}

	if err := t.Execute(out, in); err != nil {
		return errors.Wrapf(err, "execute (%T, in: %s)", t, in)
	}

	p.dependencies = t.Dependencies()

	return nil
}
