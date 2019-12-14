package handler

import (
	"filestore/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

			// 获取token中的用户信息
			token := r.Header.Get("Authorization")
			claims, err := util.ParseToken(token, []byte(pwdSalt))
			if nil != err || token == "" {
				resp := util.RespMsg{Msg: "这种行为需要登录才能使用，你的 Token 有问题!!!"}
				util.Logerr(w.Write(resp.JSONBytes()))
				return
			}
			uid := claims.(jwt.MapClaims)["uid"].(float64)
			username := claims.(jwt.MapClaims)["username"].(string)
			//fmt.Println(claims, uid)

			r.Form.Set("uid", fmt.Sprintf("%.0f", uid))
			r.Form.Set("username", username)

			handlerFunc(w, r)
		},
	)
}
