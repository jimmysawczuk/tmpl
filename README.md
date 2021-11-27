# tmpl

[![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/tmpl)](https://goreportcard.com/report/github.com/jimmysawczuk/tmpl)

**tmpl** is a small command-line utility to execute Go templates (defined in [`text/template`](https://golang.org/pkg/text/template) and [`html/template`](https://golang.org/pkg/html/template)) on files.

## Features

By default, tmpl processes configuration from a file named `tmpl.config.json` in the working directory. You can override this file path with the `-f` flag.

Here's a sample configuration file:

```jsonc
[
	{
		"in": "index.tmpl",
		"out": "out/index.html",
		"format": "html", // Can be "html", "json", or "".
		"options": {
			// Whether or not to minify the output (default: false; has no effect
			// when the format is "")
			"minify": true,

			// Environment variables to pass in for use in the template. You can
			// also set environment variables normally; variables set in the config
			// file take precedence.
			"env": {
				"FOO": "BAR"
			}
		}
	}
]
```

## Functions

In addition to the [built-in functions](https://pkg.go.dev/text/template#hdr-Functions) provided by the `text/template` package, these functions are available in every template:

| Function       |
| -------------- |
| `add`          |
| `asset`        |
| `autoreload`   |
| `env`          |
| `file`         |
| `formatTime`   |
| `getJSON`      |
| `inline`       |
| `jsonify`      |
| `markdown`     |
| `now`          |
| `parseTime`    |
| `ref`          |
| `safeCSS`      |
| `safeHTML`     |
| `safeHTMLAttr` |
| `safeJS`       |
| `seq`          |
| `sub`          |
| `svg`          |
| `timeIn`       |

## Watch mode

Watch mode (`-w`) watches all of the templates in your config for changes and rebuilds them when they're changed. Additionally, any files referenced in your templates via `ref` or similar template functions will trigger a rebuild of the template.

## Server mode

Server mode (`-s`) is the same as watch mode except it also spins up a webserver that will serve the base directory.

Additionally, this webserver has an endpoint (`/__tmpl`) which will resolve when a change is made. You can use the `autoreload` function in your template to automatically reload the page when this endpoint resolves.

## Subcommand

You can pass in a subcommand to be run by providing the `--` flag and then your command. You might want to use this if you need to run a second development process, like webpack, alongside your templates.

Here's an example:

```sh
$ tmpl -w -- webpack -w --mode=development
```

## License

[MIT](/LICENSE)
