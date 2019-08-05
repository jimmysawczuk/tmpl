package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var config = struct {
	watch   bool
	out     string
	fmt     string
	envFile string

	timestampAssets bool
}{
	timestampAssets: true,
}

func init() {
	flag.BoolVar(&config.watch, "w", false, "run continuously, watching for changes and rewriting file as needed")
	flag.StringVar(&config.out, "out", "", "output destination ('' means stdout)")
	flag.StringVar(&config.fmt, "fmt", "", "output format ('', 'html' or 'json')")
	flag.StringVar(&config.envFile, "env-file", "", "pipe .env file before executing template")
	flag.BoolVar(&config.timestampAssets, "timestamp-assets", true, "set to false to not automatically timestamp production assets")
}

func main() {
	o, err := newPayload()
	if err != nil {
		fatalErr(errors.Wrap(err, "build payload"))
	}

	flag.Parse()

	if config.envFile != "" {
		if err := godotenv.Load(config.envFile); err != nil {
			fatalErr(errors.Wrapf(err, "load .env file %s", config.envFile))
		}
	}

	var ins []io.Reader
	var prefixes []string
	if flag.Arg(0) == "--" {
		ins = append(ins, os.Stdin)
	} else {
		matches, err := doublestar.Glob(flag.Arg(0))
		if err != nil {
			fatalErr(errors.Wrapf(err, "glob: %s", flag.Arg(0)))
		}

		for _, m := range matches {
			fi, err := os.Stat(m)
			if err != nil {
				log.Println(errors.Wrapf(err, "stat file %s", m))
				continue
			}

			if fi.IsDir() {
				continue
			}

			fp, err := os.Open(m)
			if err != nil {
				log.Println(errors.Wrapf(err, "open file %s", m))
				continue
			}

			ins = append(ins, fp)
			prefixes = append(prefixes, fp.Name())
		}
	}

	prefix := longestCommonPrefix(prefixes)
	log.Println("prefix", prefix)

	var outs []io.Writer
	for _, v := range ins {
		var out io.Writer = os.Stdout

		if config.out != "" {
			if len(ins) > 1 {
				fi, err := os.Stat(config.out)
				if err != nil {
					if os.IsNotExist(err) {
						if err := os.MkdirAll(config.out, 0755); err != nil {
							fatalErr(errors.Wrapf(err, "create directory %s", config.out))
						}

						fi, err = os.Stat(config.out)
						if err != nil {
							fatalErr(errors.Wrapf(err, "stat output %s", config.out))
						}
					} else {
						fatalErr(errors.Wrapf(err, "stat output %s", config.out))
					}
				}

				if !fi.IsDir() {
					fatalErr(errors.New("directory is required for more than one input"))
				}

				name := strings.Replace(v.(*os.File).Name(), prefix, "", 1)

				outpath := filepath.Clean(config.out + "/" + name)

				out, err = os.OpenFile(outpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					fatalErr(err)
				}
			}
		}

		outs = append(outs, out)
	}

	for i := range ins {
		switch config.fmt {
		case "html":
			if err := writeHTML(ins[i], o, outs[i]); err != nil {
				fatalErr(err)
			}

		case "json":
			if err := writeJSON(ins[i], o, outs[i]); err != nil {
				fatalErr(err)
			}

		default:
			if err := writeText(ins[i], o, outs[i]); err != nil {
				fatalErr(err)
			}
		}
	}

	// var out io.Writer = os.Stdout
	// if config.output != "" {
	// 	fp, err := openOutputFile(config.output)
	// 	if err != nil {
	// 		fatalErr(errors.Wrapf(err, "open output file %s", config.output))
	// 	} else {
	// 		out = fp
	// 	}

	// 	out = fp
	// }

	// switch config.fmt {
	// case "html":
	// 	if err := writeHTML(string(by), o, out); err != nil {
	// 		fatalErr(err)
	// 	}
	// case "json":
	// 	if err := writeJSON(string(by), o, out); err != nil {
	// 		fatalErr(err)
	// 	}
	// default:
	// 	if err := writeText(string(by), o, out); err != nil {
	// 		fatalErr(err)
	// 	}
	// }
}

func openOutputFile(path string) (io.Writer, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, errors.Wrap(err, "mkdir")
	}

	fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok {
			return nil, errors.Wrap(err, "open file")
		}

		return nil, perr
	}

	return fp, nil
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(2)
}

func longestCommonPrefix(ins []string) string {
	prefix := ""
	j := 0
	for {
		for i := 1; i < len(ins); i++ {
			if j > len(ins[0]) || j > len(ins[i]) {
				return prefix
			}

			f0 := ins[0]
			fi := ins[i]
			if f0[j] != fi[j] {
				return prefix
			}
		}

		prefix = prefix + string(ins[0][j])
		j++
	}
}
