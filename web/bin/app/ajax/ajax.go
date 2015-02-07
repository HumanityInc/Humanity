package ajax

import (
	"../../db"
	"../../model"
	"../../session"
	"../../utils"
	"strings"
)

const (
	PASSWORD_MIN_LENGTH = 1

	TRIM = " \r\n\t"
)

type (
	Result struct {
		Res   int         `json:"res"`
		Error string      `json:"error,omitempty"`
		Data  interface{} `json:"data,omitempty"`
	}
)

var (
	errEmailAlreadyExists     = "EMAIL_ALREADY_EXISTS"
	errInvalidEmailAddress    = "INVALID_EMAIL"
	errInvalidPassword        = "INVALID_PASSWORD"
	errInvalidEmailOrPassword = "INVALID_EMAIL_OR_PASSWORD"
	errPasswordsNotEqual      = "PASSWORDS_NOT_EQUAL"
	errNotAuth                = "NOTAUTH"
)

func Register(c *model.Client) {

	email := c.Req.FormValue("email")
	password := c.Req.FormValue("password")
	password2 := c.Req.FormValue("password2")
	lastName := c.Req.FormValue("last_name")
	firstName := c.Req.FormValue("first_name")
	ip := c.Ip()

	if utils.IsEmail(email) {

		if len(password) >= PASSWORD_MIN_LENGTH {

			if password == password2 {

				lastName = strings.Trim(lastName, TRIM)
				firstName = strings.Trim(firstName, TRIM)

				userId, err := db.RegisterUser(email, firstName, lastName, password, ip, true)
				if err == nil {

					if user, err := db.GetUserById(userId); err == nil {
						ukey := session.SetUserCookie(c.Res)
						session.SetUser(*user, ukey)
					}

					c.WriteJson(&Result{})

				} else {
					c.WriteJson(&Result{Res: 1, Error: errEmailAlreadyExists})
				}

			} else {
				c.WriteJson(&Result{Res: 1, Error: errPasswordsNotEqual})
			}

		} else {
			c.WriteJson(&Result{Res: 1, Error: errInvalidPassword})
		}

	} else {
		c.WriteJson(&Result{Res: 1, Error: errInvalidEmailAddress})
	}
}

func Login(c *model.Client) {

	email := c.Req.FormValue("email")
	password := c.Req.FormValue("password")

	if utils.IsEmail(email) {

		if password != "" {

			user, err := db.GetUserByEmail(email, password, true)
			if err == nil {

				ukey := session.SetUserCookie(c.Res)
				session.SetUser(*user, ukey)

				c.WriteJson(&Result{})

			} else {
				c.WriteJson(&Result{Res: 1, Error: errInvalidEmailOrPassword})
			}

		} else {
			c.WriteJson(&Result{Res: 1, Error: errInvalidPassword})
		}

	} else {
		c.WriteJson(&Result{Res: 1, Error: errInvalidEmailAddress})
	}
}

func SetEmail(c *model.Client) {

	if c.User == nil {
		c.WriteJson(&Result{Res: 1, Error: errNotAuth})
		return
	}

	email := c.Req.FormValue("email")

	if utils.IsEmail(email) {

		user, err := db.SetUserEmail(c.User.Id, email)
		if err == nil {
			ukey := session.SetUserCookie(c.Res)
			session.SetUser(*user, ukey)
		}

		c.WriteJson(&Result{})

	} else {
		c.WriteJson(&Result{Res: 1, Error: errInvalidEmailAddress})
	}
}

func Whoami(c *model.Client) {

	c.WriteJson(c.User)
}

func Logout(c *model.Client) {

	session.Delete(c.Req, c.Res)

	c.WriteJson(&Result{})
}
