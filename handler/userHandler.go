package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
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
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	token := r.Header.Get("Authorization")

	// 2. 验证token是否有效
	claims, err := util.ParseToken(token, []byte(pwdSalt))
	if nil != err {
		fmt.Println(" err :", err)
		return
	}
	username := claims.(jwt.MapClaims)["username"].(string)

	// 3. 查询用户信息
	umeta, err := db.GetUserInfoDB(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed get user meta , err : " + err.Error())
		return
	}

	//// 验证uid
	uid := claims.(jwt.MapClaims)["uid"]
	if uid != umeta.Uid {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.RespMsg{Msg: "Token 错误!!!"}
		util.Logerr(w.Write(resp.JSONBytes()))
		return
	}

	// 4. 组装并响应用户数据
	w.WriteHeader(http.StatusOK)

	resp := util.RespMsg{
		Msg:  "这是搞什么事情啊???",
		Data: umeta,
	}

	util.Logerr(w.Write(resp.JSONBytes()))
}
