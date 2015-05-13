package app

// TODO use https://github.com/vkholodkov/nginx-upload-module

import (
	"../config"
	"../model"
	"../paypal"
	"../session"
	"../template_cache"
	"./ajax"
	"./oauth"
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

	cb_fn func(cb *model.Client)
)

var page_callback map[string]cb_fn
var cache = template_cache.CreateCache(TEMPLATE_PATH)
var mc *memcache.Client

func init() {

	page_callback = map[string]cb_fn{
		``:          Index,
		`index_ng`:  FeedNg,
		`grid`:      FeedNg,
		`crowdfund`: FeedNg,
		`touch`:     FeedNg,

		// `feed`:             Feed,
		// `crowdfund`:        Crowdfund,
		`create`:           Create,
		`auth`:             oauth.Auth,
		`j_login`:          ajax.Login,
		`j_logout`:         ajax.Logout,
		`j_register`:       ajax.Register,
		`j_setemail`:       ajax.SetEmail,
		`j_whoami`:         ajax.Whoami,
		`j_feed`:           ajax.Feed,
		`j_crowdfund`:      ajax.CrowdfundSave,
		`test_menu`:        TestMenu,
		`test_scroll`:      TestScroll,
		`feed_ng`:          FeedNg,
		`j_reset`:          ajax.ResetPasswd,
		`j_resetlink`:      ajax.SendResetLink,
		`j_crowdfund_info`: ajax.CrowdfundInfo,
		`j_avatar`:         ajax.SaveAvatar,
		`j_favorit`:        ajax.Favorit,
		`j_search`:         ajax.Search,
		`j_username`:       ajax.SaveUserName,
		`reset`:            ResetPasswd,
		`paypal`:           paypal.InstantPaymentNotification,
	}

	conf := config.GetConfig()

	mc = memcache.New(conf.Storage.Memcache...)

	if err := mc.Set(&memcache.Item{Key: "test", Value: []byte("test")}); err != nil {
		fmt.Println("memcache", err)
		os.Exit(1)
	}

	session.Init(mc)
}

func Bind(base string, cb_fn func(cb *model.Client)) (ok bool) {
	if page_callback != nil {
		page_callback[base] = cb_fn
		ok = true
	}
	return
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

	cb := model.Client{
		Start: time.Now().UnixNano(),
		Path:  strings.Split(strings.Trim(req.URL.Path, `/`), `/`),
		Res:   res,
		Req:   req,
	}

	defer error_check(&cb)

	if user, ok := session.GetUser(req); ok {
		cb.User = &user
	}

	if len(cb.Path) == 0 {

		Index(&cb)

	} else {

		if fn, ok := page_callback[cb.Path[0]]; ok {
			fn(&cb)
		} else {
			NotFound(&cb)
		}
	}

	cb.Stop = time.Now().UnixNano()
	fmt.Printf("%.03f %v\n", float32(cb.Stop-cb.Start)/1000000.0, cb.Path)
}

func error_check(cb *model.Client) {

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

		fmt.Println(err_str)

		InternalServerError(cb)

		// TODO send alert

		var buf [10240]byte
		runtime.Stack(buf[:], false)
		fmt.Printf("%s\n", buf)
	}
}
