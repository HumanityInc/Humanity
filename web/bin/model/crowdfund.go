package model

import (
	"fmt"
	"net/url"
	"strings"
)

const ()

type (
	Crowdfund struct {
		Id         int64    `json:"id"`
		OwnerId    int64    `json:"owner"`
		CreateTime int64    `json:"ctime"`
		UpdateTime int64    `json:"utime"`
		Goal       int64    `json:"goal"`
		Сollected  int64    `json:"сollected"`
		Visible    int16    `json:"visible"`
		Name       string   `json:"name"`
		Cover      string   `json:"cover"`
		Images     []string `json:"images"`
	}
)

func init() {

	return

	// get youtube id

	info, _ := url.Parse("https://www.youtube.com/watch?v=PnlafTs7nos&list=PLAJjZKeHBIoSsmwRSkNOhAI8D36zH8ta9")
	// info, _ := url.Parse("http://youtu.be/PnlafTs7nos?list=PLAJjZKeHBIoSsmwRSkNOhAI8D36zH8ta9")

	values := info.Query()

	switch info.Host {
	case "www.youtube.com", "youtube.com":

		id := values.Get("v")

		fmt.Println(info.Host, id)

	case "youtu.be":

		parts := strings.Split(info.RequestURI(), "?")

		if len(parts) > 0 {

			id := strings.Trim(parts[0], "/")

			fmt.Println(info.Host, id)
		}

	}

}
