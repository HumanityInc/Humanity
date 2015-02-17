package app

import (
	"../db"
	"../model"
	_ "../session"
	// "fmt"
	"html/template"
	"io/ioutil"
	// "strings"
)

func Index(c *model.Client) {

	changelog_byte, _ := ioutil.ReadFile("../htdocs/opensource/changelog.txt")
	// changelog := strings.Replace(template.HTMLEscapeString(string(changelog_byte)), "\n", "<br>", -1)

	page := Render{
		res:  c.Res,
		tmpl: "index/index.html",
		data: struct {
			User      *model.User
			Changelog template.HTML
		}{
			User:      c.User,
			Changelog: template.HTML(changelog_byte),
		},
	}
	page.Render()
}

func IndexNg(c *model.Client) {

	page := Render{
		res:  c.Res,
		tmpl: "index_ng/index.html",
		data: struct {
			User *model.User
		}{
			User: c.User,
		},
	}
	page.Render()
}

func TestMenu(c *model.Client) {

	page := Render{
		res:  c.Res,
		tmpl: "test_menu/index.html",
		data: struct {
			User *model.User
		}{
			User: c.User,
		},
	}
	page.Render()
}

func Feed(c *model.Client) {

	crowdfunds := db.Crowdfunds(0, 100)

	page := Render{
		res:  c.Res,
		tmpl: "feed/index.html",
		data: struct {
			User       *model.User
			Crowdfunds []model.Crowdfund
		}{
			User:       c.User,
			Crowdfunds: crowdfunds,
		},
	}
	page.Render()
}

func Crowdfund(c *model.Client) {

	id := c.FormValueInt64("id")

	if crowdfund, err := db.GetCrowdfund(id); err == nil {

		embedLink := crowdfund.EmbedLink()

		page := Render{
			res:  c.Res,
			tmpl: "crowdfund/index.html",
			data: struct {
				User      *model.User
				Crowdfund *model.Crowdfund
				EmbedLink string
			}{
				User:      c.User,
				Crowdfund: crowdfund,
				EmbedLink: embedLink,
			},
		}
		page.Render()

		return
	}

	NotFound(c)
}

func Create(c *model.Client) {

	id := c.FormValueInt64("id")

	var crowdfund model.Crowdfund

	if cf, err := db.GetCrowdfund(id); err == nil {
		crowdfund = *cf
	}

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
			Crowdfund:   crowdfund,
		},
	}
	page.Render()
}
