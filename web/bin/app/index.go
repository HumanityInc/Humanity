package app

import (
	"../db"
	"../model"
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

func Feed(c *model.Client) {

	page := Render{
		res:  c.Res,
		tmpl: "feed/index.html",
		data: struct {
			User *model.User
		}{
			User: c.User,
		},
	}
	page.Render()
}

func Crowdfund(c *model.Client) {

	id := c.FormValueInt64("id")

	if crowdfund, err := db.GetCrowdfund(id); err == nil {

		page := Render{
			res:  c.Res,
			tmpl: "crowdfund/index.html",
			data: struct {
				User      *model.User
				Crowdfund *model.Crowdfund
			}{
				User:      c.User,
				Crowdfund: crowdfund,
			},
		}
		page.Render()
	}

	c.Redirect("/")
}

func Create(c *model.Client) {

	id := c.FormValueInt64("id")

	page := Render{
		res:  c.Res,
		tmpl: "create/index.html",
		data: struct {
			User        *model.User
			Crowdfund   model.Crowdfund
			CrowdfundId int64
		}{
			User:        c.User,
			CrowdfundId: id,
		},
	}
	page.Render()
}
