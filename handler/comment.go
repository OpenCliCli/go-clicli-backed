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

func GetPostComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	if comments, err := db.GetPostComments(pid); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		util.SendSuccess(w, comments)
	}
}

func CommentPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))

	req, _ := ioutil.ReadAll(r.Body)
	pbody := &def.Comment{}
	if err := json.Unmarshal(req, pbody); err != nil {
		util.SendMsg(w, 400, "param error")
		return
	}

	resp, err := db.CreateComment(util.DecodeTokenID(r), pid, pbody.Content)
	if err != nil {
		util.SendMsg(w, 500, "create comment error"+err.Error())
		return
	} else {
		util.SendSuccess(w, resp)
	}
}
