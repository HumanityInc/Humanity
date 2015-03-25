package ajax

import (
	"../../db"
	"../../model"
	"../../session"
	_ "encoding/json"
	"fmt"
)

const ()

var (
	errParseJson = "PARSE_IMAGES_JSON"
	errNotSaved  = "NOT_SAVED"
	errNotFound  = "NOT_FOUND"
)

func Feed(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	offset := c.FormValueUint("offset")
	limit := c.FormValueUint("limit")
	prof_id := c.FormValueInt64("prof_id")

	if limit == 0 {
		limit = 60
	}

	if limit > 100 {
		limit = 100
	}

	if prof_id == 0 {

		list := db.Crowdfunds(c.User.Id, offset, limit)
		c.WriteJson(&Result{Data: list})

	} else {

		list := db.UserCrowdfunds(prof_id, offset, limit)
		c.WriteJson(&Result{Data: list})
	}
}

func SaveUserName(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	firstName := c.Req.FormValue("first_name")
	lastName := c.Req.FormValue("last_name")
	location := c.Req.FormValue("location")

	password := c.Req.FormValue("password")
	passwordNew1 := c.Req.FormValue("password_new1")
	passwordNew2 := c.Req.FormValue("password_new2")

	c.User.FirstName = firstName
	c.User.LastName = lastName
	c.User.Location = location

	ukey := session.GetSession(c.Req)
	session.SetUser(*c.User, ukey)

	db.SetUserName(c.User.Id, firstName, lastName)
	db.SetUserLocation(c.User.Id, location)

	if password != "" {

		if db.CheckUserPassword(c.User.Id, password) {

			if passwordNew1 == passwordNew2 {

				if db.UpdateUserPassword(c.User.Id, passwordNew1) {

					//

				} else {

					// error
				}

			} else {

				// error
			}

		} else {

			// error
		}
	}

	c.WriteJson(&Result{})
}

func SaveAvatar(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	if avatar := c.Req.FormValue("avatar"); avatar != "" {

		c.User.Picture = avatar

		ukey := session.GetSession(c.Req)
		session.SetUser(*c.User, ukey)

		db.SetUserPicture(c.User.Id, avatar)
		c.WriteJson(&Result{})

	} else {
		c.WriteJson(&Result{Res: 1})
	}
}

func CrowdfundInfo(c *model.Client) {

	id := c.FormValueInt64("id")

	if crowdfund, err := db.GetCrowdfund(id); err == nil {

		crowdfund.Video = crowdfund.EmbedLink()

		c.WriteJson(&Result{Data: crowdfund})

	} else {
		c.WriteJson(&Result{Res: 1, Error: errNotFound})
	}
}

func Follow(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	id := c.FormValueInt64("id")
	flag := false

	if c.Req.FormValue("flag") == "true" {
		flag = true
	}

	db.Follow(c.User.Id, id, flag)

	c.WriteJson(&Result{})
}

func Favorit(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	id := c.FormValueInt64("id")
	flag := false

	if c.Req.FormValue("flag") == "true" {
		flag = true
	}

	db.Favorit(c.User.Id, id, flag)

	c.WriteJson(&Result{})
}

func Search(c *model.Client) {

	q := c.Req.FormValue("q")
	list := []model.User{}

	if q != "" {
		list = db.Search(q)
	}

	c.WriteJson(&Result{Data: list})
}

func CrowdfundSave(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	images := []string{}

	// data_images := c.Req.FormValue("images")

	// if err := json.Unmarshal([]byte(data_images), &images); err != nil {
	// 	c.WriteJson(&Result{Res: 1, Error: errParseJson})
	// 	return
	// }

	crowdfund := model.Crowdfund{
		Id:      c.FormValueInt64("id"),
		OwnerId: c.User.Id,
		Goal:    c.FormValueInt64("goal"),
		Name:    c.Req.FormValue("name"),
		Video:   c.Req.FormValue("video"),
		Cover:   c.Req.FormValue("cover"),
		Images:  images,
	}

	fmt.Println(crowdfund)

	if err := db.SaveCrowdfund(&crowdfund); err == nil {

		c.WriteJson(&Result{Data: crowdfund})

	} else {
		c.WriteJson(&Result{Res: 1, Error: errNotSaved})
	}
}
