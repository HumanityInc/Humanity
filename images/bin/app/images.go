package app

import (
	"fmt"
	"github.com/disintegration/imaging"
	"os"
)

const (
	IMAGE_SOURCE = ""
	IMAGE_POSTER = "01"
)

type Size struct {
	Width  int
	Height int
}

var resize = map[string]Size{
	IMAGE_POSTER: Size{Width: 640, Height: 360},
}

// TODO use devicePixelRatio

func (c *Callback) wImages() {

	//  0      1  2  3  4
	// [images 1a 24 XX 1a244753a91626afbfdfed2c03f315ff.jpg]

	if len(c.path) == 5 {

		path := FILE_PATH + c.path[1] + "/" + c.path[2]

		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {

				key := c.path[3]

				if size, ok := resize[key]; ok {

					file := path + "/" + c.path[4]

					if info, err = os.Stat(file); err == nil {

						path += "/" + key
						file_dst := path + "/" + c.path[4]

						if info, err = os.Stat(path); err != nil {
							os.MkdirAll(path, os.ModePerm)
						}

						img, err := imaging.Open(file)
						if err != nil {
							c.InternalServerError(err.Error())
							return
						}

						img = imaging.Thumbnail(img, size.Width, size.Height, imaging.Lanczos)

						err = imaging.Save(img, file_dst)
						if err != nil {
							c.InternalServerError(err.Error())
							return
						}

						c.redirect(fmt.Sprintf("/%s/%s/%s/%s/%s", c.path[0], c.path[1], c.path[2], c.path[3], c.path[4]))
						return
					}
				}
			}
		}
	}

	c.NotFound()
}
