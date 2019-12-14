package main

import (
	"filestore/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	TestCreateToken()
}

func TestCreateToken() {
	token, _ := util.CreateToken([]byte("SecretKey"), "YDQ", 2222, "niracler")
	fmt.Println(token)

	claims, err := util.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjgsIndlYiI6Imdyb3VwdGVuIiwiZXhwIjoxNTc0OTI5NTQzLCJpYXQiOjE1NzQ5MjU5NDN9.cXxMKTKIlAJcjx9GBrjJVA_Ie-a1zW9dEBhicmgrxS8", []byte("onlinemusic"))
	if nil != err {
		fmt.Println(" err :", err)
	}
	fmt.Println("claims:", claims)
	fmt.Println("claims uid:", claims.(jwt.MapClaims)["uid"])
}
