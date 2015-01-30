package app

import (
	"../session"
	"fmt"
)

var (
	jsonMimeType = `application/json`
	jsonSuccess  = []byte(`{"res":0}`)
)

func (c *Client) jRegister() {

	email := c.req.FormValue("email")
	password := c.req.FormValue("password")
	first_name := c.req.FormValue("first_name")
	last_name := c.req.FormValue("last_name")
	ip := c.ip()

	fmt.Println(email, password, first_name, last_name, ip)
}

func (c *Client) jLogin() {

	email := c.req.FormValue("email")
	password := c.req.FormValue("password")

	fmt.Println(email, password)
}

func (c *Client) jLogout() {

	session.Delete(c.req, c.res)

	c.res.Header().Add(`Content-Type`, `application/json`)
	c.res.Write(jsonSuccess)
}
