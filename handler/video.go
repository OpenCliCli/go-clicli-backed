package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cliclitv/go-clicli/db"
	"github.com/cliclitv/go-clicli/def"
	"github.com/cliclitv/go-clicli/util"
	"github.com/julienschmidt/httprouter"
)

func AddVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	body := &def.Video{}

	if err := json.Unmarshal(req, body); err != nil {
		util.SendMsg(w, 400, "参数解析失败")
		return
	}

	if resp, err := db.AddVideo(body.Oid, body.Title, body.Content, body.Pid, body.Uid); err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendSuccess(w, resp)
	}

}

func UpdateVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	vid, _ := strconv.Atoi(id)
	req, _ := ioutil.ReadAll(r.Body)
	body := &def.Video{}

	if err := json.Unmarshal(req, body); err != nil {
		util.SendMsg(w, 400, "参数解析失败")
		return
	}

	if resp, err := db.UpdateVideo(vid, body.Oid, body.Title, body.Content, body.Pid, body.Uid); err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendSuccess(w, resp)
	}

}

func GetVideos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pid, _ := strconv.Atoi(r.URL.Query().Get("pid"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	resp, _ := db.GetVideos(pid, page, pageSize)
	util.SendSuccess(w, resp)
}

func GetVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, _ := strconv.Atoi(p.ByName("id"))
	resp, err := db.GetVideo(vid)
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendSuccess(w, resp)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	pid, _ := strconv.Atoi(r.URL.Query().Get("pid"))

	err := db.DeleteVideo(id, pid)
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendMsg(w, 200, "删除成功")
	}
}
