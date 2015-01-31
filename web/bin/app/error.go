package app

import (
	"../model"
	"net/http"
)

func MethodNotAllowed(c *model.Client) {

	c.Res.WriteHeader(http.StatusMethodNotAllowed)
	page := Render{res: c.Res, tmpl: "errors/405.html"}
	page.Render()
}

func ServiceUnavailable(c *model.Client) {

	c.Res.WriteHeader(http.StatusServiceUnavailable)
	page := Render{res: c.Res, tmpl: "errors/503.html"}
	page.Render()
}

func NotFound(c *model.Client) {

	c.Res.WriteHeader(http.StatusNotFound)
	page := Render{res: c.Res, tmpl: "errors/404.html"}
	page.Render()
}

func InternalServerError(c *model.Client) {

	c.Res.WriteHeader(http.StatusInternalServerError)
	page := Render{res: c.Res, tmpl: "errors/500.html"}
	page.Render()
}
