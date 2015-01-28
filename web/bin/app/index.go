package app

import (
	_ "../session"
)

func (c *Client) wIndex() {

	page := Render{
		res:  c.res,
		tmpl: "index/index.html",
	}
	page.Render()
}

func (c *Client) wFeed() {

	page := Render{
		res:  c.res,
		tmpl: "feed/index.html",
	}
	page.Render()
}
