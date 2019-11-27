package handler

import (
	"filestore-server/util"
	"fmt"
	"net/http"
)

// http请求拦截器
func HTTPInterception(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := r.ParseForm(); err != nil {
				fmt.Println(w, "ParseForm() err: "+err.Error())
				return
			}

			token := r.Header.Get("Authorization")
			if token == "" {
				resp := util.RespMsg{Msg: "Token 错误!!!"}
				util.Logerr(w.Write(resp.JSONBytes()))
				return
			}

			fmt.Println(token + "我在这里动了一下!!!")

			handlerFunc(w, r)
		},
	)
}
