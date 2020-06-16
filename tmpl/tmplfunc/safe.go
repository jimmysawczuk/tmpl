package tmplfunc

import (
	"html/template"
)

// SafeHTML converts the provided string to html/template's HTML type.
// This signals to the template processor that the provided string
// does not need to be escaped.
//
// See: https://golang.org/pkg/html/template#HTML
func SafeHTML(s string) template.HTML { return template.HTML(s) }

// SafeAttr converts the provided string to html/template's HTMLAttr type.
// This signals to the template processor that the provided string
// does not need to be escaped.
//
// See: https://golang.org/pkg/html/template#HTMLAttr
func SafeAttr(s string) template.HTMLAttr { return template.HTMLAttr(s) }

// SafeJS converts the provided string to html/template's JS type.
// This signals to the template processor that the provided string
// does not need to be escaped.
//
// See: https://golang.org/pkg/html/template#JS
func SafeJS(s string) template.JS { return template.JS(s) }

// SafeCSS converts the provided string to html/template's CSS type.
// This signals to the template processor that the provided string
// does not need to be escaped.
//
// See: https://golang.org/pkg/html/template#CSS
func SafeCSS(s string) template.CSS { return template.CSS(s) }
