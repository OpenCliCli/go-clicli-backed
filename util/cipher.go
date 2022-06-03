package util

import (
	"crypto/md5"
	"fmt"
)

func Cipher(pwd string) string {
	tmp := md5.Sum([]byte(pwd))
	nextPwd := fmt.Sprintf("%x", tmp)
	res := md5.Sum([]byte(nextPwd + "%132yse@clicli.us+changhao2333?"))
	newPwd := fmt.Sprintf("%x", res)
	return newPwd
}
