package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

const (
	ALPHA        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ALPHA_LENGHT = len(ALPHA)
)

var (
	re_email = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))` +
		`@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
)

func init() {

	rand.Seed(time.Now().UnixNano())
}

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
