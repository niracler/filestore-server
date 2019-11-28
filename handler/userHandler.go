package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pwdSalt = "*#890"
)

// 对用户的操作
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		HTTPInterception(getInfoHandler)(w, r) //查询用户信息
	} else if r.Method == http.MethodPost {
		registerHandler(w, r) // 用户注册
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}
}

// 用户登录接口
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Println(w, "ParseForm() err: "+err.Error())
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		encPassword := util.Sha1([]byte(password + pwdSalt))

		uid, pwdChecked := db.IsValidUserDB(username, encPassword)
		if !pwdChecked {
			w.WriteHeader(http.StatusBadRequest)
			resp := util.RespMsg{Msg: "密码或用户名错误!!!"}
			util.Logerr(w.Write(resp.JSONBytes()))
		} else {
			w.WriteHeader(http.StatusOK)

			token, _ := util.CreateToken([]byte(pwdSalt), "YDQ", uid, username)
			fmt.Println(token)

			resp := util.RespMsg{
				Data: struct {
					Username string
					Token    string
				}{
					Username: username,
					Token:    token,
				},
			}

			util.Logerr(w.Write(resp.JSONBytes()))
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}
}

// 用户注册的接口
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	for key, value := range r.Form {
		fmt.Printf("获取from中的%s:%s\n", key, value)
	}

	if len(username) < 3 || len(password) < 5 {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.RespMsg{Msg: "Invalid parameter, too short!!!"}
		util.Logerr(w.Write(resp.JSONBytes()))
		return
	}

	encPassword := util.Sha1([]byte(password + pwdSalt))
	suc := db.CreateUserDB(username, encPassword)
	if suc {
		w.WriteHeader(http.StatusCreated)
		resp := util.RespMsg{Msg: "成功注册"}
		util.Logerr(w.Write(resp.JSONBytes()))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.RespMsg{Msg: "该用户可能已存在"}
		util.Logerr(w.Write(resp.JSONBytes()))
	}
}

// 获取用户信息的接口
func getInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	username := r.Form.Get("username")
	uid, err := strconv.ParseInt(r.Form.Get("uid"), 10, 64)

	var resp util.RespMsg
	// 各种错误的情况的响应
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp = util.RespMsg{Msg: "uid 参数有问题" + err.Error()}

	} else if umeta, err := db.GetUserInfoDB(username); err != nil {
		// 尝试从数据库中找出该用户
		w.WriteHeader(http.StatusBadRequest)
		resp = util.RespMsg{Msg: "找不到该用户!!!" + err.Error()}

	} else if uid != umeta.Uid {
		// 验证uid
		w.WriteHeader(http.StatusBadRequest)
		resp = util.RespMsg{Msg: "Token 错误!!!"}

	} else {
		// 组装并响应用户数据
		w.WriteHeader(http.StatusOK)
		resp = util.RespMsg{
			Msg:  "您好,这是您的用户数据!!",
			Data: umeta,
		}
	}

	util.Logerr(w.Write(resp.JSONBytes()))
}
