package db

import (
	"errors"
)

const (
	SN_GOOGLEPLUS = `gp`
	SN_FACEBOOK   = `fb`
	SN_TWITTER    = `tw`
)

type (
	User struct {
		Id         int64  `json:"id"`
		Email      string `json:"email"`
		Password   string `json:"-"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Status     string `json:"-"`
		Registered string `json:"-"`
		LastLogin  string `json:"-"`
		LastIp     string `json:"-"`
		Gender     string `json:"-"`
		Locale     string `json:"-"`
	}
)

var (
	ErrUserAlreadyExists             = errors.New(`User already exists`)
	ErrInvalidEmailAddressOrPassword = errors.New(`Invalid email address or password`)
)

func RegisterUser(email, firstName, lastName, password string) (err error) {

	return
}

func RegisterSocialUser(socialId, socialNetworkName, email, firstName, lastName string) (err error) {

	return
}

func GetUserBySocialId(socialId, socialNetworkName string) (user *User, err error) {

	return
}

func GetUserByEmail(email, password string) (user *User, err error) {

	return
}

func GetUserById(id int64) (user *User, err error) {

	return
}

func ActivateUser(key string) (err error) {

	return
}
