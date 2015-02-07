package app

import (
	"net/http"
)

func (c *Callback) MethodNotAllowed() {

	c.res.WriteHeader(http.StatusMethodNotAllowed)
	c.res.Write([]byte("405 Method Not Allowed"))
}

func (c *Callback) ServiceUnavailable() {

	c.res.WriteHeader(http.StatusServiceUnavailable)
	c.res.Write([]byte("503 Service Temporarily Unavailable"))
}

func (c *Callback) NotFound() {

	c.res.WriteHeader(http.StatusNotFound)
	c.res.Write([]byte("404 Not Found"))
}

func (c *Callback) InternalServerError(err_str string) {

	c.res.WriteHeader(http.StatusInternalServerError)
	c.res.Write([]byte("500 Runtime Error"))
}
