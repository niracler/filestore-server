package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 文件上传的接口
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
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
		meta.UpdateFileMetaDB(fileMeta)

		io.WriteString(w, "Upload success, "+fileMeta.FileName+"\nsha1:"+fileMeta.FileSha1)

		//http.Redirect(w, r, "/file/suc/", http.StatusFound)
	}
}

// 文件上传成功的接口
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload success, hahaha")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

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

// 下载文件的接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileSha1+"\"")
	w.Write(data)
}

// 更新文件元信息的接口
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 文件删除的接口
func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)

	io.WriteString(w, "Delete success, hahaha")
}
