package app

import (
	"../config"
	"../db"
	"../session"
	"../template_cache"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	TEMPLATE_PATH = `../template/`
)

type (
	Render struct {
		res  http.ResponseWriter
		data interface{}
		tmpl string
		lang string
	}

	Client struct {
		start int64
		stop  int64
		path  []string
		res   http.ResponseWriter
		req   *http.Request
		user  *db.User
	}

	cb_fn func(cb *Client)
)

var page_callback map[string]cb_fn
var cache = template_cache.CreateCache(TEMPLATE_PATH)
var mc *memcache.Client

func init() {

	page_callback = map[string]cb_fn{
		``:           func(cb *Client) { cb.wIndex() },
		`feed`:       func(cb *Client) { cb.wFeed() },
		`auth`:       func(cb *Client) { cb.wAuth() },
		`j_login`:    func(cb *Client) { cb.jLogin() },
		`j_logout`:   func(cb *Client) { cb.jLogout() },
		`j_register`: func(cb *Client) { cb.jRegister() },
	}

	conf := config.GetConfig()

	mc = memcache.New(conf.Storage.Memcache...)

	if err := mc.Set(&memcache.Item{Key: "test", Value: []byte("test")}); err != nil {
		fmt.Println("memcache", err)
		os.Exit(1)
	}

	session.Init(mc)
}

func (render *Render) Render() {

	if render.lang == "" {
		render.lang = "en"
	}

	tmpl_path := render.lang + "/" + render.tmpl

	tmpl, _ := cache.Get(tmpl_path)

	if tmpl == nil {
		fmt.Fprint(render.res, tmpl_path+" not found")
		return
	}

	err := tmpl.Execute(render.res, render.data)
	if err != nil {
		fmt.Println(err)
	}
}

func Handler(res http.ResponseWriter, req *http.Request) {

	cb := Client{
		start: time.Now().UnixNano(),
		path:  strings.Split(strings.Trim(req.URL.Path, `/`), `/`),
		res:   res,
		req:   req,
	}

	defer error_check(&cb)

	if user, ok := session.GetUser(req); ok {
		cb.user = &user
	}

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

func error_check(cb *Client) {

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

		// TODO send alert, write log

		var buf [10240]byte
		runtime.Stack(buf[:], false)
		fmt.Printf("%s\n", buf)
	}
}
