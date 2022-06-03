package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cliclitv/go-clicli/db"
	"github.com/cliclitv/go-clicli/def"
	"github.com/cliclitv/go-clicli/util"
	"github.com/julienschmidt/httprouter"
	"github.com/nilslice/jwt"
)

func register(name string, pwd string) (err string, res bool) {
	if !util.CheckUserName(name) {
		return "用户名只能为数字字母下划线，且长度在6-16位之间", false
	}

	if !util.CheckPassword(pwd) {
		return "密码6-16位数字字母下划./", false
	}

	resp, _ := db.GetUser(name, 0)
	if resp != nil {
		return "用户名已被占用", false
	}

	if err := db.CreateSimpleUser(name, pwd); err != nil {
		return err.Error(), false
	} else {
		return "注册成功", true
	}
}

func login(name string, password string, resp *def.User) (*def.LoginSuccess, error) {
	status, _ := strconv.Atoi(resp.Status)

	if status == 0 && resp.Level > 0 {
		return nil, errors.New("用户状态异常，禁止登录")
	}

	pwd := util.Cipher(password)

	if pwd != resp.Pwd {
		return nil, errors.New("账号或密码错误")
	}

	claims := map[string]interface{}{
		"exp":   time.Now().Add(time.Hour * 12 * 7).Unix(),
		"level": resp.Level,
		"name":  resp.Name,
		"uid":   resp.Id,
	}

	token, err := jwt.New(claims)
	if err != nil {
		return nil, err
	}

	resp.Pwd = ""

	return &def.LoginSuccess{Token: token, User: resp}, nil
}

func LoginOrRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	ubody := &def.User{}

	if err := json.Unmarshal(req, ubody); err != nil {
		util.SendMsg(w, 400, err.Error())
		return
	}

	resp, _ := db.GetUserByName(ubody.Name)
	if resp == nil {
		msg, res := register(ubody.Name, ubody.Pwd)
		if !res {
			util.SendMsg(w, 400, msg)
			return
		}
	}

	loginRes, err := login(ubody.Name, ubody.Pwd, resp)
	if loginRes != nil {
		util.SendSuccess(w, loginRes)
	} else {
		if err != nil {
			util.SendMsg(w, 400, err.Error())
			return
		}
	}
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	ubody := &def.User{}

	if err := json.Unmarshal(req, ubody); err != nil {
		util.SendMsg(w, 400, "参数解析失败")
		return
	}

	if msg, res := register(ubody.Name, ubody.Pwd); res {
		util.SendMsg(w, 200, msg)
	} else {
		util.SendMsg(w, 400, msg)
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	ubody := &def.User{}

	if err := json.Unmarshal(req, ubody); err != nil {
		util.SendMsg(w, 400, err.Error())
		return
	}

	resp, err := db.GetUserByName(ubody.Name)
	if err != nil {
		util.SendMsg(w, 400, err.Error())
		return
	}

	if resp == nil {
		util.SendMsg(w, 400, "账号或密码错误")
		return
	}

	res, err := login(ubody.Name, ubody.Pwd, resp)
	if err != nil {
		util.SendMsg(w, 400, err.Error())
		return
	}

	domain := "upv.life"
	if os.Getenv("ENV") == "development" {
		domain = "127.0.0.1"
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   res.Token,
		Domain:  domain,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 7),
		MaxAge:  int((time.Hour * 24 * 7).Seconds()),
	})

	util.SendSuccess(w, res)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pint, _ := strconv.Atoi(p.ByName("id"))
	user, _ := db.GetUser("", pint)

	if !util.DouLevelChecker(w, r, user.Level, user.Id) {
		return
	}

	req, _ := ioutil.ReadAll(r.Body)
	ubody := &def.User{}
	if err := json.Unmarshal(req, ubody); err != nil {
		util.SendStandardMsg(w, http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	var realLevel int
	token := r.Header.Get(util.AUTHORIZATION_KEY)
	s := jwt.GetClaims(token)
	l := int(s["level"].(float64))
	id := int(s["uid"].(float64))

	if ubody.Level > user.Level && id == user.Id {
		util.SendMsg(w, 403, "不可自我革职")
		return
	}

	if ubody.Level < l {
		util.SendMsg(w, 403, "垂直越权")
		return
	}

	if user.Status != ubody.Status && id == pint {
		util.SendMsg(w, 403, "越权")
		return
	}

	if ubody.Pwd != "" && !util.CheckPassword(ubody.Pwd) {
		util.SendMsg(w, http.StatusBadRequest, "密码6-16位数字字母下划./")
		return
	}

	realLevel = ubody.Level
	resp, _ := db.UpdateUser(pint, ubody.Nickname, ubody.Pwd, realLevel, ubody.QQ, ubody.Bio, ubody.Status, ubody.Avatar)
	util.SendSuccess(w, resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid, _ := strconv.Atoi(p.ByName("id"))
	user, _ := db.GetUser("", uid)

	if uid == user.Id {
		util.SendMsg(w, 403, "权限不够")
		return
	}
	if !util.DouLevelChecker(w, r, user.Level, user.Id) {
		return
	}

	err := db.DeleteUser(uid)
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendMsg(w, 200, "封禁成功")
	}
}

func DeleteUserByIds(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	ids := make([]interface{}, 0)
	json.Unmarshal(req, &ids)

	if db.DeleteUserByIds(ids) == nil {
		util.SendMsg(w, 200, "done")
	} else {
		util.SendMsg(w, 500, "")
	}
}

func GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid, _ := strconv.Atoi(p.ByName("id"))
	if user, err := db.GetUser("", uid); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		util.SendSuccess(w, user)
	}
}

func GetUserByIdOrName(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := r.URL.Query().Get("uname")
	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	if user, err := db.GetUser(uname, uid); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		if user == nil {
			util.SendMsg(w, 404, "用户不存在")
			return
		}
		util.SendSuccess(w, user)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	level, _ := strconv.Atoi(r.URL.Query().Get("level"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	order := r.URL.Query().Get("order")
	name := r.URL.Query().Get("name")
	status := r.URL.Query().Get("status")

	if pageSize > 100 {
		pageSize = 15
		return
	}

	resp, _ := db.GetUsers(name, level, page, pageSize, order, status)
	util.SendSuccess(w, resp)
}

func SearchUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.URL.Query().Get("key")

	resp, err := db.SearchUsers(key)
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		res := &def.Users{Users: resp}
		util.SendSuccess(w, res)
	}
}

func GetUserStat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid, _ := strconv.Atoi(p.ByName("id"))
	if user, err := db.GetUserStat(uid); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		util.SendSuccess(w, user)
	}
}

func GetMe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token := r.Header.Get(util.AUTHORIZATION_KEY)
	s := jwt.GetClaims(token)
	if s["uid"] == nil {
		util.SendMsg(w, 403, "unauthorized")
		return
	}
	uid := int(s["uid"].(float64))
	if user, err := db.GetUser("", uid); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		util.SendSuccess(w, user)
	}
}
