package app

import (
	"../db"
	"../model"
	_ "../session"
	"../utils"
	_ "fmt"
	"html/template"
	"io/ioutil"
	"strings"
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
			ProfId    uint64
			Code      string
		}{
			User:      c.User,
			Changelog: template.HTML(changelog_byte),
		},
	}
	page.Render()
}

func ResetPasswd(c *model.Client) {

	if len(c.Path) == 3 {

		changelog_byte, _ := ioutil.ReadFile("../htdocs/opensource/changelog.txt")

		prof_id := utils.ToUint64(c.Path[1])
		code := c.Path[2]

		page := Render{
			res:  c.Res,
			tmpl: "index/index.html",
			data: struct {
				User      *model.User
				ProfId    uint64
				Code      string
				Changelog template.HTML
			}{
				User:      c.User,
				ProfId:    prof_id,
				Code:      code,
				Changelog: template.HTML(changelog_byte),
			},
		}
		page.Render()

		return
	}

	c.Redirect("/")
}

func FeedNg(c *model.Client) {

	if c.User == nil {
		c.Redirect("/")
		return
	}

	if db.Invitee[strings.ToLower(c.User.Email)] == 1 {

		page := Render{
			res:  c.Res,
			tmpl: "feed_ng/index.html",
			data: struct {
				User *model.User
			}{
				User: c.User,
			},
		}
		page.Render()
		return
	}

	c.Redirect("/")
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

func TestScroll(c *model.Client) {

	page := Render{
		res:  c.Res,
		tmpl: "test_scroll/index.html",
		data: struct {
			User *model.User
		}{
			User: c.User,
		},
	}
	page.Render()
}

func Feed(c *model.Client) {

	// if c.User == nil {
	// 	c.WriteJson(&Result{Res: 1, Error: errNotAuth})
	// 	return
	// }

	// crowdfunds := db.Crowdfunds(0, 100)

	// page := Render{
	// 	res:  c.Res,
	// 	tmpl: "feed/index.html",
	// 	data: struct {
	// 		User       *model.User
	// 		Crowdfunds []model.Crowdfund
	// 	}{
	// 		User:       c.User,
	// 		Crowdfunds: crowdfunds,
	// 	},
	// }
	// page.Render()
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
