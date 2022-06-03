package util

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	auth "github.com/nilslice/jwt"
)

var AUTHORIZATION_KEY = "Authorization"

type AuthInfo struct {
	Uid   int
	Level int
}

func DecodeAuthorization(r *http.Request) *AuthInfo {
	token := r.Header.Get(AUTHORIZATION_KEY)
	if token == "" {
		return &AuthInfo{
			Uid:   -1,
			Level: -1,
		}
	}
	s := auth.GetClaims(token)

	return &AuthInfo{
		Uid:   int(s["uid"].(float64)),
		Level: int(s["level"].(float64)),
	}
}

func DecodeTokenID(r *http.Request) int {
	token := r.Header.Get(AUTHORIZATION_KEY)
	if token == "" {
		return -1
	}
	s := auth.GetClaims(token)

	return int(s["uid"].(float64))
}

func DecodeTokenLevel(r *http.Request) int {
	token := r.Header.Get(AUTHORIZATION_KEY)

	if token == "" {
		return -1
	}

	s := auth.GetClaims(token)

	return int(s["level"].(float64))
}

/*
	level: 持有者的权限
	target: 目标权限
*/
func LevelChecker(level int, target int) bool {
	return level <= target
}

// OwnerChecker 判断内容所属 不判断 root 和 admin eg: admin 可改 root 文章
func OwnerChecker(w http.ResponseWriter, r *http.Request, id int) bool {
	log.Println("ownerChecker", id)
	token := r.Header.Get(AUTHORIZATION_KEY)
	s := auth.GetClaims(token)
	uid := int(s["uid"].(float64))
	level := int(s["level"].(float64))

	if level <= 1 { // root 为所欲为
		return true
	}
	if uid != id {
		SendMsg(w, 403, "越权操作")
	}

	return uid == id
}

// DouLevelChecker root VS admin -> 判断内容所属 & level eg: admin 不可改 root 资料 （需要取得 uid 的 level）
func DouLevelChecker(w http.ResponseWriter, r *http.Request, authLevel int, uid int) bool {
	log.Println("douLevelChecker", authLevel, uid)
	token := r.Header.Get(AUTHORIZATION_KEY)
	s := auth.GetClaims(token)
	id := int(s["uid"].(float64))
	level := int(s["level"].(float64))

	// root 为所欲为
	if level == 0 {
		return true
	} else {
		if uid == id {
			return true
		}

		if !LevelChecker(level, authLevel-1) {
			SendMsg(w, 403, "同级越权")
			return false
		}
		return true
	}
}

/******* router middleware *******/

//   0: '超级管理员',
//   1: '管理员',
//   2: '创作者',
//   3: '注册用户',
//   10: '全部',

func LevelAuth(h httprouter.Handle, level int) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get(AUTHORIZATION_KEY)

		if token != "" && auth.Passes(token) {
			claims := auth.GetClaims(token)
			_level := int(claims["level"].(float64))

			if _level <= level {
				h(w, r, ps)
				return
			}

			SendMsg(w, 403, http.StatusText(http.StatusUnauthorized))
			return
		}

		SendMsg(w, 401, http.StatusText(http.StatusUnauthorized))
	}
}

func Authorization(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get(AUTHORIZATION_KEY)

		if token == "" || !auth.Passes(token) {
			SendMsg(w, 401, http.StatusText(http.StatusUnauthorized))
			return
		}

		h(w, r, ps)
	}
}

func PostFilter(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		if DecodeTokenLevel(r) > 1 {
			for i := range ps {
				if ps[i].Key == "status" {
					ps[i].Value = "3"
				}
			}
		}

	}
}
