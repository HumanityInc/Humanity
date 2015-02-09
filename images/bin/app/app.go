package app

// TODO use https://github.com/vkholodkov/nginx-upload-module

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DOMAIN   = ` TODO `
	BASE_URL = `https:// TODO /`
)

type (
	Callback struct {
		start int64
		stop  int64
		path  []string
		res   http.ResponseWriter
		req   *http.Request
	}

	cb_fn func(cb *Callback)
)

var page_callback map[string]cb_fn

func init() {

	page_callback = map[string]cb_fn{
		``:       func(cb *Callback) { cb.wIndex() },
		`upload`: func(cb *Callback) { cb.wUpload() },
		`images`: func(cb *Callback) { cb.wImages() },
	}
}

func Handler(res http.ResponseWriter, req *http.Request) {

	cb := Callback{
		start: time.Now().UnixNano(),
		path:  strings.Split(strings.Trim(req.URL.Path, `/`), `/`),
		res:   res,
		req:   req,
	}

	defer error_check(&cb)

	if len(cb.path) == 0 {

		cb.wIndex()

	} else {

		if fn, ok := page_callback[cb.path[0]]; ok {

			fn(&cb)

		} else {

			cb.NotFound()
		}
	}

	cb.stop = time.Now().UnixNano()
	fmt.Printf("%.03f %v\n", float32(cb.stop-cb.start)/1000000.0, cb.path)
}

func error_check(cb *Callback) {

	if err := recover(); err != nil {

		pc, file, line, ok := runtime.Caller(4)

		if !ok {
			file = "?"
			line = 0
		}

		fn_name := ""
		fn := runtime.FuncForPC(pc)

		if fn == nil {
			fn_name = "?()"
		} else {
			dot_name := filepath.Ext(fn.Name())
			fn_name = strings.TrimLeft(dot_name, ".") + "()"
		}

		err_str := fmt.Sprintf("%s:%d %s: %s\n", file, line, fn_name, err)

		cb.InternalServerError(err_str)
	}
}

func (c *Callback) ip() (ip string) {
	if remote := strings.Split(c.req.RemoteAddr, ":"); len(remote) > 0 {
		ip = remote[0]
	}
	return
}

func (c *Callback) redirect(uri string) {

	c.res.Header().Set("Location", uri)

	if c.req.Method == "POST" {
		c.res.WriteHeader(http.StatusSeeOther)
	} else {
		c.res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// --- --- --- --- ---
