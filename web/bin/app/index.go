package app

import (
	"../model"
	_ "../sendmail"
	_ "../session"
)

func Index(c *model.Client) {

	page := Render{
		res:  c.Res,
		tmpl: "index/index.html",
		data: struct {
			User *model.User
		}{
			User: c.User,
		},
	}
	page.Render()
}
