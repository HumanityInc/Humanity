package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	IMAGE_MAX_SIZE = 1440
	FILE_PATH      = "../htdocs/images/"
	FILE_URL       = "/images/"
)

var (
	ext_allow = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".bmp":  true,
		".tif":  true,
		".tiff": true,
	}

	mime_allow = map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/jpg":  true,
		"image/bmp":  true,
		"image/tiff": true,
	}

	ext_convert = map[string]string{
		".tif":  ".jpg",
		".tiff": ".jpg",
		".bmp":  ".jpg",
	}
)

func (c *Callback) wUpload() {

	if c.req.Method == "POST" {

		file, header, err := c.req.FormFile("file")
		if err != nil {

			c.InternalServerError(err.Error())
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))

		content_type := ""
		if file_type := header.Header["Content-Type"]; len(file_type) > 0 {
			content_type = file_type[0]
		}

		_, ext_ok := ext_allow[ext]
		_, mime_ok := mime_allow[content_type]

		if ext_ok && mime_ok {

			buf, _ := ioutil.ReadAll(file)

			hash := md5.New()

			fmt.Fprint(hash, buf)
			filename := fmt.Sprintf("%x", hash.Sum(nil))

			// TODO add slice change
			sub_path := fmt.Sprintf("%s/%s/", filename[0:2], filename[2:4])

			path := FILE_PATH + sub_path

			info, err := os.Stat(path)
			if err != nil {
				os.MkdirAll(path, os.ModePerm)
			} else if info.IsDir() {
				os.MkdirAll(path, os.ModePerm)
			}

			if ext == ".jpeg" {
				ext = ".jpg"
			}

			file_path := path + filename + ext
			file_url := FILE_URL + sub_path + filename + ext

			err = ioutil.WriteFile(file_path, buf, 0644)
			if err != nil {

				c.InternalServerError(err.Error())
				return
			}

			if to_ext, ok := ext_convert[ext]; ok {

				defer os.Remove(file_path)

				img, err := imaging.Open(file_path)
				if err != nil {

					c.InternalServerError(err.Error())
					return
				}

				file_path = path + filename + to_ext
				file_url = FILE_URL + sub_path + filename + to_ext

				bounds := img.Bounds()
				max_width, max_height, max_size := bounds.Max.X, bounds.Max.Y, IMAGE_MAX_SIZE

				if max_width > max_size || max_height > max_size {

					с := float32(max_width) / float32(max_height)

					if 1 < с {
						img = imaging.Resize(img, max_size, 0, imaging.CatmullRom)
					} else {
						img = imaging.Resize(img, 0, max_size, imaging.CatmullRom)
					}
				}

				err = imaging.Save(img, file_path)
				if err != nil {

					c.InternalServerError(err.Error())
					return
				}

			} else {

				img, err := imaging.Open(file_path)
				if err != nil {

					c.InternalServerError(err.Error())
					return
				}

				bounds := img.Bounds()
				max_width, max_height, max_size := bounds.Max.X, bounds.Max.Y, IMAGE_MAX_SIZE

				if max_width > max_size || max_height > max_size {

					c2 := float32(max_width) / float32(max_height)

					if 1 < c2 {
						img = imaging.Resize(img, max_size, 0, imaging.CatmullRom)
					} else {
						img = imaging.Resize(img, 0, max_size, imaging.CatmullRom)
					}

					err = imaging.Save(img, file_path)
					if err != nil {

						c.InternalServerError(err.Error())
						return
					}
				}
			}

			c.res.Header().Set("Content-Type", "text/plain; charset=utf-8")

			link := map[string]string{
				"filelink": file_url,
			}

			data, _ := json.Marshal(&link)
			c.res.Write(data)
		}

	} else {

		c.MethodNotAllowed()
	}
}
