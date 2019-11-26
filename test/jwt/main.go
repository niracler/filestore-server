package main

import (
	"filestore-server/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	TestCreateToken()
}

func TestCreateToken() {
	token, _ := util.CreateToken([]byte("SecretKey"), "YDQ", 2222, "niracler")
	fmt.Println(token)

	claims, err := util.ParseToken(token, []byte("SecretKey"))
	if nil != err {
		fmt.Println(" err :", err)
	}
	fmt.Println("claims:", claims)
	fmt.Println("claims uid:", claims.(jwt.MapClaims)["uid"])
}
