package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jimmysawczuk/tmpl/config"
	"github.com/jimmysawczuk/tmpl/pipe"
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

func init() {
	flag.Usage = func() {
		fmt.Printf("tmpl %s; built %s (rev. %s)\n\n", version, date, revision)

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
	fp, err := os.Open(configFile)
	if err != nil {
		return errors.Wrapf(err, "open config file (path: %s)", configFile)
	}

	var blocks []config.Block
	if err := json.NewDecoder(fp).Decode(&blocks); err != nil {
		return errors.Wrap(err, "json: decode config file")
	}

	watcher, err := pipe.New(watchMode)
	if err != nil {
		return errors.Wrap(err, "watch: new")
	}
	defer watcher.Close()

	mode := tmpl.ModeProduction
	if watchMode {
		mode = tmpl.ModeLocal
	}

	pipes := []*pipe.Pipe{}

	for i, b := range blocks {

		pipe := &pipe.Pipe{
			Format: b.Format,
			Mode:   mode,

			Minify: b.Options.Minify,
			Env:    b.Options.Env,
		}

		pipe.In, _ = filepath.Abs(b.In)

		ostat, err := os.Stat(b.Out)
		if err == nil {
			if ostat.IsDir() {
				outpath := filepath.Join(b.Out, filepath.Base(b.In))

				_, err := os.Stat(outpath)
				if err == nil || os.IsNotExist(err) {
					pipe.Out, _ = filepath.Abs(outpath)
				} else {
					return errors.Wrapf(err, "stat output path (path: %s, block: %d)", outpath, i+1)
				}

				pipe.Out, _ = filepath.Abs(filepath.Join(b.Out, filepath.Base(b.In)))
			} else {
				pipe.Out, _ = filepath.Abs(b.Out)
			}
		} else if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(b.Out), 0755); err != nil {
				return errors.Wrapf(err, "mkdir (path: %s, block: %d)", filepath.Dir(b.Out), i+1)
			}

			pipe.Out, _ = filepath.Abs(b.Out)
		} else {
			return errors.Wrapf(err, "stat (output: %s, block: %d)", b.Out, i+1)
		}

		if err := watcher.AddPipe(pipe); err != nil {
			log.Printf("couldn't watch path %s: %s", pipe.In, err)
		}

		pipes = append(pipes, pipe)
	}

	for _, pipe := range pipes {
		if err := pipe.Run(); err != nil {
			return errors.Wrapf(err, "run pipeline (path: %s)", pipe.In)
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
		go startServer(port, watcherCh)
	}

	if watchMode {
		done := make(chan bool)
		go watcher.Watch(watcherCh)
		<-done
	} else if cmd != nil {
		err := cmd.Wait()
		if err != nil {
			return errors.Wrapf(err, "command: wait (%s)", strings.Join(runCommand, " "))
		}
	}

	return nil
}

func startServer(port int, watcherCh chan string) {
	log.Printf("starting server on http://localhost:%d", port)

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

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
