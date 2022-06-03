package util

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func SendSuccess(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(200)
	resStr, _ := json.Marshal(Response{Code: 200, Data: res})
	io.WriteString(w, string(resStr))
}

func SendStandardMsg(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	resStr, _ := json.Marshal(Response{Code: code, Msg: http.StatusText(code)})
	io.WriteString(w, string(resStr))
}

func SendMsg(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	resStr, _ := json.Marshal(Response{Code: code, Msg: msg})
	io.WriteString(w, string(resStr))
}
