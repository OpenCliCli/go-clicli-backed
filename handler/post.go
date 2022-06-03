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

func GetPostById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	uid := util.DecodeTokenID(r)
	resp, _ := db.GetPostByIdWithUser(pid, uid)
	util.SendSuccess(w, resp)
}

func SearchPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.URL.Query().Get("key")
	status := r.URL.Query().Get("status")
	if util.DecodeTokenLevel(r) == -1 {
		status = "3"
	}
	resp, _ := db.FindPostByTitleOrContent(key, status)

	util.SendSuccess(w, resp)
}

func GetPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	tag := r.URL.Query().Get("tag")
	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	key := r.URL.Query().Get("key")
	t := r.URL.Query().Get("type")
	auth := util.DecodeAuthorization(r)

	if pageSize == 0 || pageSize > 30 {
		pageSize = 15
	}

	if page == 0 {
		page = 1
	}

	if auth.Level == -1 {
		status = "3"
	}

	resp, err := db.GetPosts(page, pageSize, status, sort, tag, uid, key, t, auth.Uid)
	if err != nil {
		util.SendMsg(w, 500, err.Error())
		return
	}

	util.SendSuccess(w, resp)
}

func AddPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	pbody := &def.Post{}

	if err := json.Unmarshal(req, pbody); err != nil {
		util.SendMsg(w, 400, "param error")
		return
	}

	authInfo := util.DecodeAuthorization(r)
	if pbody.Type == "video" && authInfo.Level > 2 {
		util.SendMsg(w, 403, "denied")
		return
	}

	postStatus := pbody.Status
	if authInfo.Level > 1 {
		if pbody.Type == "video" {
			postStatus = "2"
		} else {
			postStatus = "3"
		}
	}

	resp, err := db.CreatePost(pbody.Title, pbody.Content, postStatus, pbody.Sort, pbody.Tag, authInfo.Uid, pbody.Type)
	if _, err := db.AddVideos(pbody.Videos, authInfo.Uid, resp.Id); err != nil {
		util.SendMsg(w, 500, "video add failed")
		return
	}

	if err != nil {
		util.SendMsg(w, 500, "")
		return
	} else {
		util.SendSuccess(w, resp)
	}

}

func UpdatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid := p.ByName("id")
	pint, _ := strconv.Atoi(pid)
	post, _ := db.GetPostById(pint)

	if !util.OwnerChecker(w, r, post.Uid) {
		return
	}

	req, _ := ioutil.ReadAll(r.Body)
	pbody := &def.Post{}
	if err := json.Unmarshal(req, pbody); err != nil {
		util.SendMsg(w, 400, "参数解析失败")
		return
	}

	if util.DecodeTokenLevel(r) > 1 {
		pbody.Status = post.Status
	}
	resp, err := db.UpdatePost(pint, pbody.Title, pbody.Content, pbody.Status, pbody.Sort, pbody.Tag, pbody.Type)
	if pbody.Type == "video" {
		if _, err := db.AddVideos(pbody.Videos, post.Uid, resp.Id); err != nil {
			util.SendMsg(w, 500, "视频添加失败")
			return
		}
	}
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendSuccess(w, resp)
	}

}

func DeletePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	post, _ := db.GetPostById(pid)
	if !util.OwnerChecker(w, r, post.Uid) {
		return
	}

	err := db.DeletePost(pid)
	if err != nil {
		util.SendMsg(w, 500, "数据库错误")
		return
	} else {
		util.SendMsg(w, 200, "已下架")
	}
}

func DeletePostByIds(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, _ := ioutil.ReadAll(r.Body)
	ids := make([]interface{}, 0)
	json.Unmarshal(req, &ids)

	if db.DeletePostByIds(ids) == nil {
		util.SendMsg(w, 200, "done")
	} else {
		util.SendMsg(w, 500, "")
	}
}

func GetRank(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if resp, err := db.GetRankList(); err != nil {
		util.SendMsg(w, 500, err.Error())
	} else {
		util.SendSuccess(w, resp)
	}
}

func FindPostByTag(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tag := r.URL.Query().Get("tag")
	type_ := r.URL.Query().Get("type")
	if type_ == "" {
		type_ = "video"
	}

	posts, _ := db.FindPostByTag(tag, type_)

	util.SendSuccess(w, posts)
}

//点赞
func LikePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	post, _ := db.GetPostById(pid)
	if post == nil {
		util.SendMsg(w, 404, "post not found")
		return
	}
	uid := util.DecodeTokenID(r)

	if err := db.LikePost(pid, uid); err == nil {
		util.SendMsg(w, 200, "done")
	} else {
		util.SendMsg(w, 500, err.Error())
	}
}

//取消点赞
func UnLikePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	uid := util.DecodeTokenID(r)

	if db.UnLikePost(pid, uid) == nil {
		util.SendMsg(w, 200, "done")
	} else {
		util.SendMsg(w, 500, "")
	}
}

//收藏文章
func CollectPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	post, _ := db.GetPostById(pid)
	if post == nil {
		util.SendMsg(w, 404, "post not found")
		return
	}

	uid := util.DecodeTokenID(r)
	if err := db.CollectPost(pid, uid); err == nil {
		util.SendMsg(w, 200, "done")
	} else {
		util.SendMsg(w, 400, err.Error())
	}
}

func UnCollectPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	uid := util.DecodeTokenID(r)
	if err := db.UnCollectPost(pid, uid); err != nil {
		util.SendMsg(w, 400, err.Error())
	} else {
		util.SendMsg(w, 200, "done")
	}
}

//是否收藏
func IsCollectedPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pid, _ := strconv.Atoi(p.ByName("id"))
	uid := util.DecodeTokenID(r)
	if isCollected, err := db.IsCollectedPost(pid, uid); err != nil {
		util.SendMsg(w, 400, err.Error())
	} else {
		util.SendSuccess(w, isCollected)
	}
}

//我的收藏
func GetUserCollections(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid, _ := strconv.Atoi(p.ByName("id"))
	posts, err := db.FindUserCollectPost(uid)
	if err != nil {
		util.SendMsg(w, 500, "")
		return
	}
	util.SendSuccess(w, posts)
}
