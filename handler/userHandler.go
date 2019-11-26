package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	pwdSalt = "*#890"
)

// 对用户的操作
func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, "你已经成功获取用户信息了")
		return
	} else if r.Method == http.MethodPost { // 用户注册
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
			w.Write([]byte("Invalid parameter, too short!!!"))
			return
		}

		encPassword := util.Sha1([]byte(password + pwdSalt))
		suc := db.UserSignup(username, encPassword)
		if suc {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("成功注册"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("该用户可能已存在"))
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}

}

// 登录接口
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Println(w, "ParseForm() err: "+err.Error())
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		encPassword := util.Sha1([]byte(password + pwdSalt))

		pwdChecked := db.UserSignin(username, encPassword)
		if !pwdChecked {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("密码或用户名错误!!!"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(genToken(username)))
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}
}

func genToken(username string) string {
	// md5(username+timestamp+tokenSalt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + pwdSalt))
	return tokenPrefix + ts[:8]
}
