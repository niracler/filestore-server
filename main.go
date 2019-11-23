package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/", handler.UploadHandler)
	http.HandleFunc("/file/suc/", handler.UploadSucHandler)
	http.HandleFunc("/file/meta/", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download/", handler.DownloadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server , err: %s", err.Error())
	}
	fmt.Println("Hello World")
}
