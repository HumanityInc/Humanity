package ajax

import (
	"../../db"
	"../../model"
	_ "encoding/json"
	"fmt"
)

const ()

var (
	errParseJson = "PARSE_IMAGES_JSON"
	errNotSaved  = "NOT_SAVED"
)

func Feed(c *model.Client) {

	offset := c.FormValueUint("offset")
	limit := c.FormValueUint("limit")

	if limit == 0 {
		limit = 60
	}

	if limit > 100 {
		limit = 100
	}

	list := db.Crowdfunds(offset, limit)

	c.WriteJson(Result{Data: list})
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
