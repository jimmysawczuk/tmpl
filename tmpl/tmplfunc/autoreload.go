package tmplfunc

import "html/template"

// Autoreload returns an HTML snippet that you can embed in your templates to
// automatically reload the page when a change is detected.
func Autoreload(m Moder) func() template.HTML {
	return func() template.HTML {
		if m.IsProduction() {
			return template.HTML("")
		}

		return SafeHTML(`<script>fetch('/__tmpl').then(function(){console.log("Change detected, reloading!");top.location.reload()})</script>`)
	}
}
