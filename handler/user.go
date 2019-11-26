package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
)

const (
	pwdSalt = "*#890"
)

// 对用户的操作
func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, "你已经成功获取用户信息了")
		return
	}

	// 用户注册
	if r.Method == http.MethodPost {
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
	}

}
