package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cliclitv/go-clicli/handler"
	"github.com/cliclitv/go-clicli/util"
	"github.com/julienschmidt/httprouter"
	"github.com/nilslice/jwt"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENV") == "development" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "false")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,upv-user-agent")
	}

	m.r.ServeHTTP(w, r)
}

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/me", util.Authorization(handler.GetMe))
	router.POST("/user/register", handler.Register)
	router.POST("/user/login", handler.Login)
	router.POST("/user/login/or/register", handler.LoginOrRegister)
	router.POST("/user/update/:id", util.Authorization(handler.UpdateUser))
	router.POST("/user/delete/:id", util.LevelAuth(handler.DeleteUser, 1))
	router.POST("/user/deleteUserByIds", util.LevelAuth(handler.DeleteUserByIds, 1))
	router.GET("/users", handler.GetUsers)
	router.GET("/user/:id", handler.GetUser)
	router.GET("/user/:id/stat", handler.GetUserStat)
	router.GET("/user", handler.GetUserByIdOrName)
	router.GET("/user/:id/collections", handler.GetUserCollections)

	/** 文章 **/
	router.GET("/rank", handler.GetRank)
	router.GET("/posts", handler.GetPosts)
	router.GET("/post/:id", handler.GetPostById)
	router.GET("/search/posts", handler.SearchPosts)
	router.GET("/posts/recommends", handler.FindPostByTag)
	router.POST("/post/add", util.Authorization(handler.AddPost))
	router.GET("/post/:id/comments", handler.GetPostComments)
	router.POST("/post/comment/:id", util.Authorization(handler.CommentPost))

	router.POST("/post/isCollected/:id", util.Authorization(handler.IsCollectedPost))
	router.POST("/post/like/:id", util.Authorization(handler.LikePost))
	router.POST("/post/unlike/:id", util.Authorization(handler.UnLikePost))
	router.POST("/post/collect/:id", util.Authorization(handler.CollectPost))
	router.POST("/post/uncollect/:id", util.Authorization(handler.UnCollectPost))
	router.GET("/posts/recommend", handler.FindPostByTag)
	router.POST("/post/delete/:id", util.LevelAuth(handler.DeletePost, 2))
	router.POST("/user/deletePostByIds", util.LevelAuth(handler.DeletePostByIds, 2))
	router.POST("/post/update/:id", util.LevelAuth(handler.UpdatePost, 2))

	router.POST("/video/add", util.LevelAuth(handler.AddVideo, 2))
	router.POST("/video/update/:id", util.LevelAuth(handler.UpdateVideo, 2))
	router.POST("/video/delete", util.LevelAuth(handler.DeleteVideo, 2))
	router.GET("/video/:id", handler.GetVideo)
	router.GET("/videos", handler.GetVideos)
	router.GET("/search/users", handler.SearchUsers)
	router.GET("/pv/:pid", handler.GetPv)

	router.GET("/jx", handler.Jx)

	router.POST("/upload", util.LevelAuth(handler.UploadFile, 2))
	if os.Getenv("ENV") == "development" {
		router.ServeFiles("/upload/*filepath", http.Dir("upload"))
	}

	return router
}

func main() {
	str := util.RandStr(10)
	jwt.Secret([]byte(str))

	mh := NewMiddleWareHandler(RegisterHandler())
	port := os.Getenv("BACKEND_PORT")
	log.Fatal(http.ListenAndServe(":"+port, mh))
}
