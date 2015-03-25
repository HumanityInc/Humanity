package model

import ()

const (
	SN_FACEBOOK   = 1
	SN_TWITTER    = 2
	SN_GOOGLEPLUS = 3
)

type (
	User struct {
		Id         int64  `json:"id"`
		Email      string `json:"email"`
		Password   string `json:"-"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Location   string `json:"location"`
		Picture    string `json:"picture"`
		Activate   string `json:"-"`
		Registered string `json:"-"`
		LastLogin  string `json:"-"`
		LastIp     string `json:"-"`
		Gender     string `json:"-"`
	}

	SocialProfile struct {
		Id        string `json:"id"`
		SnId      int    `json:"sn_id"`
		UserId    int64  `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   string `json:"picture"`
		Link      string `json:"link"`
		Gender    string `json:"gender"`
		LastIp    string `json:"-"`
	}
)
