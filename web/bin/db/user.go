package db

import (
	"../model"
	"../sendmail"
	"../utils"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const ()

var (
	ErrUserAlreadyExists             = errors.New(`User already exists`)
	ErrUserNotFound                  = errors.New(`User not found`)
	ErrInvalidEmailAddressOrPassword = errors.New(`Invalid email address or password`)

	Invitee = map[string]int{}
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

func CheckUserPassword(userId int64, password string) (ok bool) {

	hashedPassword := ""

	err := db.QueryRow(`SELECT "password" FROM "public"."profiles" WHERE "id"=?`, userId).Scan(&hashedPassword)

	if err != nil {
		logger.Println(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return
	}

	ok = true
	return
}

func UpdateUserPassword(userId int64, password string) (ok bool) {

	hash := hashPassword(password)

	_, err := db.Exec(`UPDATE "public"."profiles" SET "password"=$1 WHERE "id"=$2`,
		hash, userId)

	if err != nil {
		logger.Println(err)
		return
	}

	ok = true
	return
}

//

func Search(q, ip string) (list []model.User) {

	opts := &sphinx.Options{
		Host:       "127.0.0.1",
		Port:       9312,
		MatchMode:  sphinx.SPH_MATCH_ANY,
		Timeout:    500,
		MaxMatches: 1000,
	}

	userIdsInt64 := []int64{}
	userIds := []string{}

	if sphinxClient := sphinx.NewClient(opts); sphinxClient != nil {

		defer sphinxClient.Close()

		sphinxClient.SetLimits(0, 24, 300, 0)

		if result, err := sphinxClient.Query(q, "src1", ip); err == nil {

			fmt.Println(result.Matches, err)

			resultCount := len(result.Matches)

			userIds = make([]string, 0, resultCount)

			for _, match := range result.Matches {

				userIdsInt64 = append(userIdsInt64, int64(match.DocId))
				userIds = append(userIds, fmt.Sprintf("%d", match.DocId))
			}

		} else {

			fmt.Println("[SEARCH]", err)
			userIdsInt64 = []int64{8, 11, 13}
			userIds = []string{"8", "11", "13"}
		}
	}

	if len(userIds) == 0 {
		return
	}

	if len(userIds) > 5 {
		userIdsInt64 = userIdsInt64[:5]
		userIds = userIds[:5]
	}

	rows, err := db.Query(`SELECT ` +
		`"id", "first_name", "last_name", "last_login", "registered", "activate", "last_ip", "picture" ` +
		`FROM "public"."profiles" ` +
		`WHERE "id" IN (` + strings.Join(userIds, ",") + `) `)

	if err != nil {
		logger.Println(err)
		return
	}
	defer rows.Close()

	user := model.User{}
	list = make([]model.User, 0, 32)

	for rows.Next() {

		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.LastLogin, &user.Registered,
			&user.Activate, &user.LastIp, &user.Picture); err != nil {

			logger.Println(err)
			return
		}

		user.FirstName = strings.Title(user.FirstName)
		user.LastName = strings.Title(user.LastName)

		list = append(list, user)
	}

	tmp := make([]model.User, 0, len(list))

	for _, id := range userIdsInt64 {
		for _, u := range list {
			if u.Id == id {
				tmp = append(tmp, u)
			}
		}
	}

	list = tmp

	return
}

func UpdatePassword(profId uint64, code, password string) (ok bool) {

	var id uint64

	err := db.QueryRow(`SELECT "prof_id" FROM "public"."reset_passwd" WHERE "prof_id"=$1 AND "code"=$2 LIMIT 1`,
		profId, code).Scan(&id)

	if err != nil {
		logger.Println(err)
		return
	}

	if id != 0 {

		hash := hashPassword(password)

		_, err := db.Exec(`UPDATE "public"."profiles" SET "password"=$1 WHERE "id"=$2`,
			hash, id)

		if err != nil {
			logger.Println(err)
			return
		}

		_, err = db.Exec(`DELETE FROM "public"."reset_passwd" WHERE "prof_id"=$1`, id)

		if err != nil {
			logger.Println(err)
			return
		}

		ok = true
	}
	return
}

func SendResetLink(email string) (ok bool) {

	if user, err := GetUserByEmail(email, "", false); err == nil {

		code := utils.RandomString(48)

		_, err = db.Exec(`INSERT INTO "public"."reset_passwd" ("prof_id", "code") VALUES ($1, $2)`,
			user.Id, code)

		if err != nil {

			result, err := db.Exec(`UPDATE "public"."reset_passwd" SET "code"=$1 WHERE "prof_id"=$2`,
				code, user.Id)

			if err != nil {
				logger.Println(err)
				return
			}

			ar, err := result.RowsAffected()

			if err != nil {
				logger.Println(err)
				return
			}

			if err == nil && ar != 0 {

				m := sendmail.MailReset{Link: fmt.Sprintf("https://ishuman.me/reset/%d/%s", user.Id, code)}
				go m.Send(user.Email)

				ok = true
			}
			return

		} else {

			m := sendmail.MailReset{Link: fmt.Sprintf("https://ishuman.me/reset/%d/%s", user.Id, code)}
			go m.Send(user.Email)

			ok = true
		}
	}
	return
}

//

func RegisterUser(email, firstName, lastName, password, ip string, useEmail bool) (userId int64, err error) {

	hashedPassword := hashPassword(password)
	unix_time := time.Now().Unix()

	if useEmail {

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

		m := sendmail.MailSuccess{}
		go m.Send(email)

	} else {

		err = db.QueryRow(`INSERT INTO "public"."profiles" `+
			`("password", "first_name", "last_name", "last_login", "registered", "activate", "last_ip") `+
			`VALUES `+
			`($1, $2, $3, $4, $5, $6, $7) `+
			`RETURNING id`,
			hashedPassword, firstName, lastName, 0, unix_time, 0, ip).Scan(&userId)

		if err != nil {
			logger.Println(err)
			return
		}
	}

	return
}

func RegisterSocialUser(socialProfile model.SocialProfile) (user_ptr *model.User, err error) {

	var userId int64

	userId, err = RegisterUser(socialProfile.Email, socialProfile.FirstName,
		socialProfile.LastName, "", socialProfile.LastIp, socialProfile.Email != "")

	if err == nil {

		_, err = db.Exec(`INSERT INTO "public"."profiles_social" `+
			`("id", "sn_id", "user_id", "first_name", "last_name", "email", "picture", "link", "gender") `+
			`VALUES `+
			`($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			socialProfile.Id, socialProfile.SnId, userId, socialProfile.FirstName, socialProfile.LastName,
			socialProfile.Email, socialProfile.Picture, socialProfile.Link, socialProfile.Gender)

		if err != nil {
			logger.Println(err)
		}

		user_ptr, err = GetUserById(userId)

		m := sendmail.MailSuccess{}
		go m.Send(socialProfile.Email)

	} else {

		user_ptr, err = GetUserByEmail(socialProfile.Email, "", false)
	}
	return
}

//

func GetUserBySocialId(socialId string, socialNetwork int) (user_ptr *model.User, err error) {

	user := model.User{}
	email := sql.NullString{}

	fmt.Println(socialId, socialNetwork)

	err = db.QueryRow(`SELECT "id", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip", "picture", "location" `+
		`FROM "public"."profiles" `+
		`WHERE "id"=(SELECT "user_id" FROM "public"."profiles_social" WHERE "id"=$1 AND "sn_id"=$2)`,
		socialId, socialNetwork).Scan(&user.Id, &email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp, &user.Picture, &user.Location)

	if err != nil {
		logger.Println(err)
		return
	}

	if email.Valid {
		user.Email = email.String
	}

	user.FirstName = strings.Title(user.FirstName)
	user.LastName = strings.Title(user.LastName)

	user_ptr = &user
	return
}

func GetUserByEmail(email, password string, checkPassword bool) (user_ptr *model.User, err error) {

	user, hashedPassword := model.User{}, ""

	err = db.QueryRow(`SELECT `+
		`"id", "password", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip", "picture", "location" `+
		`FROM "public"."profiles" `+
		`WHERE "email"=$1`, email).Scan(&user.Id, &hashedPassword, &user.Email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp, &user.Picture, &user.Location)

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

	user.FirstName = strings.Title(user.FirstName)
	user.LastName = strings.Title(user.LastName)

	user_ptr = &user
	return
}

func SetUserLocation(id int64, location string) (err error) {

	_, err = db.Exec(`UPDATE "public"."profiles" SET "location"=$1 WHERE "id"=$2`,
		location, id)

	if err != nil {
		logger.Println(err)
		return
	}
	return
}

func SetUserName(id int64, firstName, lastName string) (err error) {

	_, err = db.Exec(`UPDATE "public"."profiles" SET "first_name"=$1, "last_name"=$2 WHERE "id"=$3`,
		firstName, lastName, id)

	if err != nil {
		logger.Println(err)
		return
	}
	return
}

func SetUserPicture(id int64, picture string) (err error) {

	_, err = db.Exec(`UPDATE "public"."profiles" SET "picture"=$1 WHERE "id"=$2`,
		picture, id)

	if err != nil {
		logger.Println(err)
		return
	}
	return
}

func SetUserEmail(id int64, email string) (user_ptr *model.User, err error) {

	result, err := db.Exec(`UPDATE "public"."profiles" SET "email"=$1 WHERE "id"=$2 AND ("email" is NULL OR "email"='')`,
		email, id)

	if err != nil {
		logger.Println(err)
		return
	}

	ar, err := result.RowsAffected()

	if err != nil {
		logger.Println(err)
		return
	}

	if err == nil && ar != 0 {

		m := sendmail.MailSuccess{}
		go m.Send(email)
	}

	user_ptr, err = GetUserById(id)
	return
}

func GetUserById(id int64) (user_ptr *model.User, err error) {

	user := model.User{}
	email := sql.NullString{}

	err = db.QueryRow(`SELECT `+
		`"id", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip", "picture", "location" `+
		`FROM "public"."profiles" `+
		`WHERE "id"=$1`, id).Scan(&user.Id, &email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp, &user.Picture, &user.Location)

	if err != nil {
		logger.Println(err)
		return
	}

	if email.Valid {
		user.Email = email.String
	}

	user.FirstName = strings.Title(user.FirstName)
	user.LastName = strings.Title(user.LastName)

	user_ptr = &user
	return
}

//

func ActivateUser(key string) (err error) {

	return
}
