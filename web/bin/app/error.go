package app

import (
	"net/http"
)

func (c *Client) MethodNotAllowed() {

	c.res.WriteHeader(http.StatusMethodNotAllowed)
	page := Render{res: c.res, tmpl: "errors/405.html"}
	page.Render()
}

func (c *Client) ServiceUnavailable() {

	c.res.WriteHeader(http.StatusServiceUnavailable)
	page := Render{res: c.res, tmpl: "errors/503.html"}
	page.Render()
}

func (c *Client) NotFound() {

	c.res.WriteHeader(http.StatusNotFound)
	page := Render{res: c.res, tmpl: "errors/404.html"}
	page.Render()
}

func (c *Client) InternalServerError(err_str string) {

	c.res.WriteHeader(http.StatusInternalServerError)
	page := Render{res: c.res, tmpl: "errors/500.html"}
	page.Render()
}
