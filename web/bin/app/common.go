package app

import (
	"net/http"
	"strconv"
	"strings"
)

//

func (c *Client) ip() (ip string) {
	if remote := strings.Split(c.req.RemoteAddr, ":"); len(remote) > 0 {
		ip = remote[0]
	}
	return
}

func (c *Client) redirect(uri string) {

	c.res.Header().Set("Location", uri)

	if c.req.Method == "POST" {
		c.res.WriteHeader(http.StatusSeeOther)
	} else {
		c.res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

//

func (c *Client) FormValueUint(name string) uint {
	return ToUint(c.req.FormValue(name))
}

func (c *Client) FormValueInt(name string) int {
	return ToInt(c.req.FormValue(name))
}

func (c *Client) FormValueUint64(name string) uint64 {
	return ToUint64(c.req.FormValue(name))
}

func (c *Client) FormValueInt64(name string) int64 {
	return ToInt64(c.req.FormValue(name))
}

//

func ToUint(val string) uint {
	u64, _ := strconv.ParseUint(val, 10, 32)
	return uint(u64)
}

func ToInt(val string) int {
	i64, _ := strconv.ParseInt(val, 10, 32)
	return int(i64)
}

func ToUint64(val string) uint64 {
	u64, _ := strconv.ParseUint(val, 10, 64)
	return u64
}

func ToInt64(val string) int64 {
	i64, _ := strconv.ParseInt(val, 10, 64)
	return i64
}
