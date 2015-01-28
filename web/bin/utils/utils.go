package utils

import (
	"math/rand"
	"regexp"
)

const (
	ALPHA        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ALPHA_LENGHT = len(ALPHA)
)

var (
	re_email = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))` +
		`@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
)

func IsEmail(email string) bool {
	return re_email.MatchString(email)
}

func RandomString(strlen int) string {
	buf := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		buf[i] = ALPHA[rand.Intn(ALPHA_LENGHT)]
	}
	return string(buf)
}
