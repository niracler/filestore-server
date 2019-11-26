package handler

import (
	"filestore-server/db"
	"filestore-server/meta"
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
	if r.Method == http.MethodGet {
		getUserInfoHandler(w, r) //查询用户信息
	} else if r.Method == http.MethodPost {
		registerHandler(w, r) // 用户注册
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}
}

// 用户登录接口
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
			util.Logerr(w.Write([]byte("密码或用户名错误!!!")))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			token, _ := util.CreateToken([]byte(pwdSalt), "YDQ", 2222, username)
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
		util.Logerr(w.Write([]byte("Invalid parameter, too short!!!")))
		return
	}

	encPassword := util.Sha1([]byte(password + pwdSalt))
	suc := db.UserSignup(username, encPassword)
	if suc {
		w.WriteHeader(http.StatusCreated)
		util.Logerr(w.Write([]byte("成功注册")))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		util.Logerr(w.Write([]byte("该用户可能已存在")))
	}
}

// 获取用户信息的接口
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("claims:", claims)
	username := claims.(jwt.MapClaims)["username"].(string)

	// 3. 查询用户信息
	umeta, err := meta.GetUserMetaDB(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed get user meta , err : " + err.Error())
		return
	}

	// 4. 组装并响应用户数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Println(umeta)

	resp := util.RespMsg{
		Msg:  "这是搞什么事情啊???",
		Data: umeta,
	}

	util.Logerr(w.Write(resp.JSONBytes()))
}
