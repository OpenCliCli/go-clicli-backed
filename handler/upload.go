package handler

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cliclitv/go-clicli/util"
	"github.com/julienschmidt/httprouter"
)

func UploadFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseMultipartForm(100 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		util.SendMsg(w, 400, "parse file error")
		return
	}
	defer file.Close()

	prefix := "./upload/" + time.Now().Format("2006/01/02") + "/"
	_, err = os.Stat(prefix)

	if os.IsNotExist(err) {
		mask := syscall.Umask(0)
		if err := os.MkdirAll(prefix, 0777); err != nil {
			util.SendMsg(w, 500, "mkdir error")
			defer syscall.Umask(mask)
			return
		}
		defer syscall.Umask(mask)
	} else if err != nil {
		util.SendMsg(w, 500, "stat error")
		return
	}

	chunk := strings.Split(handler.Filename, ".")

	f, err := os.OpenFile(prefix+chunk[0]+"_"+strconv.FormatInt(time.Now().Unix(), 10)+"."+chunk[len(chunk)-1], os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		util.SendMsg(w, 500, "open file error")
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err == nil {
		util.SendSuccess(w, f.Name())
	} else {
		util.SendMsg(w, 500, "upload error")
	}
}
