package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// 对文件操作的接口
func FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uploadHandler(w, r) // 文件上传
	} else if r.Method == http.MethodGet {
		getMetaHandler(w, r) // 获取文件元信息
	} else if r.Method == http.MethodDelete {
		deleteHandler(w, r) // 删除文件
	} else if r.Method == http.MethodPut {
		updateMetaHandler(w, r) // 更新文件元信息
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // 其他操作不允许
	}
}

// 文件上传的操作
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	//接收文件流以及存储到本地目录
	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data , err: %s\n", err.Error())
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "/tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err : %s\n", err.Error())
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed to save data into file, err: %s\n", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	meta.CreateFileMetaDB(fileMeta) // 将上传的文件的元信息更新到数据库

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload success, "+fileMeta.FileName+"\nsha1:"+fileMeta.FileSha1)
}

// 通过sha1获取文件元信息
func getMetaHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	filehash := r.Form["filehash"][0]
	//fMate := meta.GetFileMeta(filehash)
	fMate, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 文件删除的接口
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}
	fileSha1 := r.Form.Get("filehash")

	fMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusNoContent)
	io.WriteString(w, "Delete success, hahaha")
}

// 更新文件元信息的接口
func updateMetaHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	opType := r.FormValue("op")
	fileSha1 := r.FormValue("filehash")
	newFileName := r.FormValue("filename")

	for key, value := range r.Form {
		fmt.Printf("获取from中的%s:%s\n", key, value)
	}

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	curFileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curFileMeta.FileName = newFileName
	meta.UpdateFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
