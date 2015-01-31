package session

import (
	"../model"
	"../utils"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
	"strconv"
	"time"
)

const (
	MC_USER_PREFIX = "u:"
	COOKIE_NAME    = `K`
	KEY_LENGHT     = 24
	SESSION_TTL    = 7 * 60 * 60 * 24
	COOKIE_TTL     = 7
)

var mc *memcache.Client

func Init(mem *memcache.Client) {
	mc = mem
}

func SetUserCookie(res http.ResponseWriter) string {

	key := utils.RandomString(KEY_LENGHT) + strconv.FormatInt(time.Now().UnixNano(), 36)

	cookie := &http.Cookie{
		Name:     COOKIE_NAME,
		Path:     "/",
		Value:    key,
		Expires:  time.Now().Local().AddDate(0, 0, COOKIE_TTL),
		HttpOnly: true,
	}

	http.SetCookie(res, cookie)

	return key
}

func SetUser(user model.User, ukey string) (ok bool) {

	key_usr := MC_USER_PREFIX + ukey

	if data, err := json.Marshal(&user); err == nil {

		err = mc.Set(&memcache.Item{
			Key:        key_usr,
			Value:      data,
			Expiration: SESSION_TTL,
		})

		if err == nil {
			ok = true
		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}

	return
}

func getUser(ukey string) (user model.User, ok bool) {

	key_usr := MC_USER_PREFIX + ukey

	if item_url, err := mc.Get(key_usr); err == nil {

		user_str := string(item_url.Value)

		if err = json.Unmarshal([]byte(user_str), &user); err == nil {
			ok = true
		} else {
			fmt.Println(err)
		}

	} /* else {
		fmt.Println(err)
	}*/

	return
}

func GetSession(req *http.Request) (ukey string) {

	if cookie, err := req.Cookie(COOKIE_NAME); err == nil {
		ukey = cookie.Value
	}
	return
}

func GetUser(req *http.Request) (user model.User, ok bool) {

	if cookie, err := req.Cookie(COOKIE_NAME); err == nil {

		ukey := cookie.Value
		user, ok = getUser(ukey)
	}
	return
}

func Delete(req *http.Request, res http.ResponseWriter) { // logout

	if cookie, err := req.Cookie(COOKIE_NAME); err == nil {

		key_usr := MC_USER_PREFIX + cookie.Value

		mc.Delete(key_usr)

		cookie := &http.Cookie{
			Name:     COOKIE_NAME,
			Path:     "/",
			Value:    "",
			Expires:  time.Now().Local().AddDate(0, 0, -1),
			HttpOnly: true,
		}

		http.SetCookie(res, cookie)
	}
}
