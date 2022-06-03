package util

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func RandStr(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func SqlPlaceholders(n int) string {
	var b strings.Builder
	for i := 0; i < n-1; i++ {
		b.WriteString("?,")
	}
	if n > 0 {
		b.WriteString("?")
	}
	return b.String()
}

func CheckUserName(name string) bool {
	math, _ := regexp.Match(`^[a-zA-Z0-9_]{4,8}$`, []byte(name))
	return math
}

func CheckPassword(pwd string) bool {
	math, _ := regexp.Match(`^[a-zA-Z0-9./_]{6,16}$`, []byte(pwd))
	return math
}
