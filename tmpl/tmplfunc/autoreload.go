package tmplfunc

import "html/template"

// Autoreload returns an HTML snippet that you can embed in your templates to
// automatically reload the page when a change is detected.
func Autoreload() template.HTML {
	return SafeHTML(`<script>fetch('/__tmpl').then(function(){console.log("Change detected, reloading!");top.location.reload()})</script>`)
}
