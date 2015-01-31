package model

import (
	"../utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const ()

type (
	Client struct {
		Start int64
		Stop  int64
		Path  []string
		Res   http.ResponseWriter
		Req   *http.Request
		User  *User
	}
)

func (c *Client) WriteJson(v interface{}) {

	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprint(c.Res, err.Error())
		return
	}

	c.Res.Header().Add(`Content-Type`, `application/json`)
	c.Res.Write(data)
}

func (c *Client) Ip() (ip string) {
	if remote := strings.Split(c.Req.RemoteAddr, ":"); len(remote) > 0 {
		ip = remote[0]
	}
	return
}

func (c *Client) Redirect(uri string) {

	c.Res.Header().Set("Location", uri)

	if c.Req.Method == "POST" {
		c.Res.WriteHeader(http.StatusSeeOther)
	} else {
		c.Res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

//

func (c *Client) FormValueUint(name string) uint {
	return utils.ToUint(c.Req.FormValue(name))
}

func (c *Client) FormValueInt(name string) int {
	return utils.ToInt(c.Req.FormValue(name))
}

func (c *Client) FormValueUint64(name string) uint64 {
	return utils.ToUint64(c.Req.FormValue(name))
}

func (c *Client) FormValueInt64(name string) int64 {
	return utils.ToInt64(c.Req.FormValue(name))
}
