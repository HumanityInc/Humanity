package model

import ()

const ()

type (
	User struct {
		Id         int64  `json:"id"`
		Email      string `json:"email"`
		Password   string `json:"-"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Activate   string `json:"-"`
		Registered string `json:"-"`
		LastLogin  string `json:"-"`
		LastIp     string `json:"-"`
		Gender     string `json:"-"`
		Locale     string `json:"-"`
	}
)
