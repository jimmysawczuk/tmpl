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

-   [`add`](#add)
-   [`asset`](#asset)
-   [`autoreload`](#autoreload)
-   [`env`](#env)
-   [`file`](#file)
-   [`formatTime`](#formatTime)
-   [`getJSON`](#getJSON)
-   [`inline`](#inline)
-   [`jsonify`](#jsonify)
-   [`markdown`](#markdown)
-   [`now`](#now)
-   [`parseTime`](#parseTime)
-   [`ref`](#ref)
-   [`safeCSS`](#safeCSS)
-   [`safeHTML`](#safeHTML)
-   [`safeHTMLAttr`](#safeHTMLAttr)
-   [`safeJS`](#safeJS)
-   [`seq`](#seq)
-   [`sub`](#sub)
-   [`svg`](#svg)
-   [`timeIn`](#timeIn)

### `add`

> add returns the sum of the two arguments.

```
{{ add 2 2 }}
```

returns:

```
4
```

### `asset`

> asset returns the path provided. In the future, asset may gain the ability to clean or validate the path.

```
{{ asset "/css/style.css" }}
```

returns

```
/css/style.css
```

### `autoreload`

> autoreload returns an HTML snippet that you can embed in your templates to automatically reload the page when a change is detected.

```
{{ autoreload }}
```

returns

```
<script>...</script>
```

### `env`

> env returns the environment variable defined at the provided key. Variables set in `tmpl.config.json` take precedence.

```
{{ env "NODE_ENV" }}
```

returns

```
production
```

### `file`

> file loads the file at the path provided and returns its contents. It _does not_ create a ref; updating this file's contents won't trigger an update in watch mode.

```
{{ file "some-letter.txt" }}
```

returns

```
...data...
```

### `formatTime`

> formatTime formats the provided time with the provided format. You can either specify both arguments at once or pipe a `time.Time` into this function to format it.

```
{{ now | formatTime "Jan 2, 2006 3:04 PM" }}
{{ formatTime now "Jan 2, 2006 3:04 PM" }}
```

returns

```
Nov 28, 2021 10:09 AM
Nov 28, 2021 10:09 AM
```

### `getJSON`

> getJSON loads the file at the provided path and unmarshals it into a `map[string]interface{}`. It _does not_ create a ref; updating this file's contents won't trigger an update in watch mode.

```
{{ getJSON "REVISION.json" }}
```

returns

```
map[string]interface{}{
    ...
}
```

### `inline`

> inline loads the file at the path provided and returns its contents. It creates a ref so that updates to the file trigger an update in watch mode.

```
{{ file "some-letter.txt" }}
```

returns

```
...data...
```

### `jsonify`

> jsonify marshals the provided input as a JSON string.

```
{{ now | jsonify }}
```

returns

```
"2021-11-28T10:09:00Z"
```

### `markdown`

> markdown reads the file at the provided path, parses its contents as Markdown and returns the HTML. It creates a ref so that updates to the file trigger an update in watch mode.

```
{{ markdown "path-to-markdown.md" }}
```

returns

```
<h1>...</h1>

<p>...</p>
```

### `now`

> now returns the time of the template's execution in the local timezone.

```
{{ now | jsonify }}
```

returns

```
"2021-11-28T10:09:00Z"
```

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
