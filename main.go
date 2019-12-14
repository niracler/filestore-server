package main

import (
	"filestore/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/", handler.FileHandler)
	http.HandleFunc("/user/", handler.UserHandler)
	http.HandleFunc("/token/", handler.GetTokenHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server , err: %s", err.Error())
	}
}
