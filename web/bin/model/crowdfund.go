package model

import (
	"net/url"
	"strings"
)

const ()

type (
	Crowdfund struct {
		Id         int64    `json:"id"`
		NextId     int64    `json:"next"`
		PrevId     int64    `json:"prev"`
		OwnerId    int64    `json:"owner"`
		CreateTime int64    `json:"ctime"`
		UpdateTime int64    `json:"utime"`
		Goal       int64    `json:"goal"`
		Сollected  int64    `json:"сollected"`
		Visible    int16    `json:"visible"`
		Favorit    bool     `json:"favorit"`
		Name       string   `json:"name"`
		Video      string   `json:"video"`
		Cover      string   `json:"cover"`
		Images     []string `json:"images"`
	}
)

func getYoutubeId(link string) (id string) {

	info, _ := url.Parse(link)
	values := info.Query()

	switch info.Host {
	case "www.youtube.com", "youtube.com":

		if id = values.Get("v"); id == "" {

			parts := strings.Split(info.RequestURI(), "?")
			if len(parts) > 0 {

				uri := strings.Trim(parts[0], "/")
				if strings.Index(uri, "embed/") == 0 {

					id = uri[6:]
				}
			}
		}

	case "youtu.be":

		parts := strings.Split(info.RequestURI(), "?")
		if len(parts) > 0 {

			id = strings.Trim(parts[0], "/")
		}
	}

	return
}

func (cf *Crowdfund) EmbedLink() (link string) {

	if id := getYoutubeId(cf.Video); id != "" {

		param := url.Values{}
		param.Add("autoplay", "0")
		param.Add("controls", "0")
		param.Add("showinfo", "0")
		param.Add("autohide", "1")
		param.Add("rel", "0")

		link = "https://youtube.com/embed/" + id + "?" + param.Encode()
	}

	return
}

// func init() {

// 	return

// 	links := []string{
// 		"http://www.youtube.com/watch?v=xjS6SftYQaQ",
// 		"http://www.youtube.com/embed/xjS6SftYQaQ",
// 		"http://www.youtube.com/watch?v=xjS6SftYQaQ&list=SPA60DCEB33156E51F",
// 		"https://www.youtube.com/watch?v=xjS6SftYQaQ&list=PLAJjZKeHBIoSsmwRSkNOhAI8D36zH8ta9",
// 		"http://youtu.be/xjS6SftYQaQ?list=PLAJjZKeHBIoSsmwRSkNOhAI8D36zH8ta9",
// 	}

// 	for i, lnk := range links {

// 		fmt.Println(i, getYoutubeId(embedLink(getYoutubeId(lnk))))
// 	}
// }
