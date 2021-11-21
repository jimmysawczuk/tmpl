package tmplfunc

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// Markdown converts the provided Markdown to HTML.
func Markdown(in string) (string, error) {
	buf := bytes.Buffer{}

	gm := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.EnDash: []byte("&mdash;"),
				}),
			),
		),
	)

	if err := gm.Convert([]byte(in), &buf); err != nil {
		return "", errors.Wrap(err, "goldmark: convert")
	}

	return buf.String(), nil
}
