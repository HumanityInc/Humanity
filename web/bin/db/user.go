package db

import (
	"../model"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const ()

var (
	ErrUserAlreadyExists             = errors.New(`User already exists`)
	ErrUserNotFound                  = errors.New(`User not found`)
	ErrInvalidEmailAddressOrPassword = errors.New(`Invalid email address or password`)
)

func hashPassword(password string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Println(err)
		return ""
	}

	return string(hashedPassword)
}

//

func RegisterUser(email, firstName, lastName, password, ip string) (userId int64, err error) {

	hashedPassword := hashPassword(password)
	unix_time := time.Now().Unix()

	err = db.QueryRow(`INSERT INTO "public"."profiles" `+
		`("password", "first_name", "last_name", "email", "last_login", "registered", "activate", "last_ip") `+
		`VALUES `+
		`($1, $2, $3, $4, $5, $6, $7, $8) `+
		`RETURNING id`,
		hashedPassword, firstName, lastName, email, 0, unix_time, 0, ip).Scan(&userId)

	if err != nil {
		logger.Println(err)
		return
	}

	return
}

func RegisterSocialUser(socialProfile model.SocialProfile) (user_ptr *model.User, err error) {

	var userId int64

	userId, err = RegisterUser(socialProfile.Email, socialProfile.FirstName, socialProfile.LastName, "", socialProfile.LastIp)

	if err == nil {

		_, err = db.Exec(`INSERT INTO "public"."profiles_social" `+
			`("id", "sn_id", "user_id", "first_name", "last_name", "email", "picture", "link", "gender") `+
			`VALUES `+
			`($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			socialProfile.Id, socialProfile.SnId, userId, socialProfile.FirstName, socialProfile.LastName,
			socialProfile.Email, socialProfile.Picture, socialProfile.Link, socialProfile.Gender)

		if err != nil {
			logger.Println(err)
			return
		}

		user_ptr, err = GetUserById(userId)

	} else {

		user_ptr, err = GetUserByEmail(socialProfile.Email, "", false)
	}
	return
}

//

func GetUserBySocialId(socialId string, socialNetwork int) (user_ptr *model.User, err error) {

	user := model.User{}

	err = db.QueryRow(`SELECT "id", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip" `+
		`FROM "public"."profiles" `+
		`WHERE "id"=(SELECT "user_id" FROM "public"."profiles_social" WHERE "id"=$1 AND "sn_id"=$2)`,
		socialId, socialNetwork).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp)

	if err != nil {
		fmt.Println(err)
		logger.Println(err)
		return
	}

	user_ptr = &user
	return
}

func GetUserByEmail(email, password string, checkPassword bool) (user_ptr *model.User, err error) {

	user, hashedPassword := model.User{}, ""

	err = db.QueryRow(`SELECT `+
		`"id", "password", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip" `+
		`FROM "public"."profiles" `+
		`WHERE "email"=$1`, email).Scan(&user.Id, &hashedPassword, &user.Email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp)

	if err != nil {
		logger.Println(err)
		return
	}

	if checkPassword {

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			logger.Println(err)
			return
		}
	}

	user_ptr = &user
	return
}

func GetUserById(id int64) (user_ptr *model.User, err error) {

	user := model.User{}

	err = db.QueryRow(`SELECT `+
		`"id", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip" `+
		`FROM "public"."profiles" `+
		`WHERE "id"=$1`, id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp)

	if err != nil {
		logger.Println(err)
		return
	}

	user_ptr = &user
	return
}

//

func ActivateUser(key string) (err error) {

	return
}
