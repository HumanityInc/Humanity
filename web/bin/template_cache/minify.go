package template_cache

import (
	"fmt"
	"github.com/tdewolff/minify"
)

var (
	MinifyOn = true
)

func HtmlMinify(html *string) string {

	if MinifyOn {

		minify_html, err := minify.NewMinifierDefault().MinifyString("text/html", *html)
		if err != nil {

			fmt.Println("MinifyString:", err)
			return *html
		}

		return minify_html
	}

	return *html
}
