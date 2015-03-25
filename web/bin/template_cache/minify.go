package template_cache

import (
	"fmt"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	// "github.com/tdewolff/minify/trim"
)

// go get -u github.com/tdewolff/minify
// go get -u github.com/tdewolff/parse

var (
	MinifyOn = !false
)

func HtmlMinify(text *string) string {

	if MinifyOn {

		m := minify.New()

		m.AddFunc("text/html", html.Minify)
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("text/js", js.Minify)
		// m.Add("*/*", trim.Minify)

		minify_html, err := minify.String(m, "text/html", *text)
		if err != nil {
			fmt.Println("MinifyString:", err)
			return *text
		}

		return minify_html
	}

	return *text
}
