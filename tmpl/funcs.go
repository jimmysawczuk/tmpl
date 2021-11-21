package tmpl

import "github.com/jimmysawczuk/tmpl/tmpl/tmplfunc"

func (t *Tmpl) funcs() map[string]interface{} {
	return map[string]interface{}{
		"add":          tmplfunc.Add,
		"asset":        tmplfunc.Asset,
		"autoreload":   tmplfunc.Autoreload,
		"env":          tmplfunc.EnvFunc(t.envVars),
		"file":         tmplfunc.File,
		"formatTime":   tmplfunc.FormatTime,
		"getJSON":      tmplfunc.GetJSON,
		"inline":       tmplfunc.Inline(t),
		"jsonify":      tmplfunc.JSONify,
		"link":         tmplfunc.Link(t),
		"markdown":     tmplfunc.Markdown,
		"now":          tmplfunc.NowFunc(t.now),
		"parseTime":    tmplfunc.ParseTime,
		"safeCSS":      tmplfunc.SafeCSS,
		"safeHTML":     tmplfunc.SafeHTML,
		"safeHTMLAttr": tmplfunc.SafeAttr,
		"safeJS":       tmplfunc.SafeJS,
		"seq":          tmplfunc.Seq,
		"sub":          tmplfunc.Sub,
		"svg":          tmplfunc.SVG(t),
		"timeIn":       tmplfunc.TimeIn,
	}
}