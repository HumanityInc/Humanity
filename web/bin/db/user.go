package db

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// "id"         int8 NOT NULL DEFAULT nextval('profile_seq'::regclass),
// "email"      varchar(128) NOT NULL,
// "password"   varchar(40) NOT NULL DEFAULT ''::character varying,
// "first_name" varchar(64) NOT NULL DEFAULT ''::character varying,
// "last_name"  varchar(64) NOT NULL DEFAULT ''::character varying,
// "last_ip"    inet NOT NULL DEFAULT '0.0.0.0'::inet,
// "status"     int2 NOT NULL DEFAULT 0,
// "registered" int8 NOT NULL DEFAULT 0,
// "last_login" int8 NOT NULL DEFAULT 0,

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
		Activate   string `json:"-"`
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

func hashPassword(password string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Println(err)
		return ""
	}

	return string(hashedPassword)
}

//

func RegisterUser(email, firstName, lastName, password, ip string) (err error) {

	hashedPassword := hashPassword(password)
	unix_time := time.Now().Unix()

	_, err = db.Exec(`INSERT INTO "public"."profiles" `+
		`("password", "first_name", "last_name", "email", "last_login", "registered", "activate", "last_ip") `+
		`VALUES `+
		`($1, $2, $3, $4, $5, $6, $7, $8)`,
		hashedPassword, firstName, lastName, email, 0, unix_time, 0, ip)

	if err != nil {
		logger.Println(err)
		return
	}

	return
}

func RegisterSocialUser(socialId, socialNetworkName, email, firstName, lastName string) (err error) {

	return
}

//

func GetUserBySocialId(socialId, socialNetworkName string) (user *User, err error) {

	return
}

func GetUserByEmail(email, password string) (user_ptr *User, err error) {

	user, hashedPassword := User{}, ""

	err = db.QueryRow(`SELECT `+
		`"id", "password", "email", "first_name", "last_name", "last_login", "registered", "activate", "last_ip" `+
		`FROM "public"."profiles" `+
		`WHERE "email"=$1`, email).Scan(&user.Id, &hashedPassword, &user.Email, &user.FirstName, &user.LastName,
		&user.LastLogin, &user.Registered, &user.Activate, &user.LastIp)

	if err != nil {
		logger.Println(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		logger.Println(err)
		return
	}

	user_ptr = &user
	return
}

func GetUserById(id int64) (user_ptr *User, err error) {

	user := User{}

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
