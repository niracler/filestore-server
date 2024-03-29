package handler

import (
	"filestore/db"
	"filestore/meta"
	"filestore/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 对文件操作的接口
func FileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		HTTPInterception(uploadHandler)(w, r) // 文件上传
	} else if r.Method == http.MethodGet {
		HTTPInterception(QueryFileHandler)(w, r) // 获取文件元信息
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

	// 接收文件流
	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data , err: %s\n", err.Error())
		return
	}
	defer file.Close()

	// 保存文件信息
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "/tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 存储到文件目录
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

	// 获取用户名
	username := r.Form.Get("username")

	//更新文件表
	msg := ""
	isExist, err1 := db.CreateFileDB(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
	if isExist {
		msg = msg + "伪秒传,"
	}

	//更新用户文件表
	isExist, err2 := db.CreateUserFileDB(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	var resp util.RespMsg
	if err1 == nil && err2 == nil && !isExist {
		w.WriteHeader(http.StatusCreated)
		resp = util.RespMsg{
			Msg:  msg + "上传文件成功",
			Data: fileMeta,
		}
	} else if isExist {
		w.WriteHeader(http.StatusBadRequest)
		resp = util.RespMsg{Msg: "该文件已存在"}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		resp = util.RespMsg{Msg: "出现了什么问题, 失败了"}
	}

	util.Logerr(w.Write(resp.JSONBytes()))
}

// 通过token获取用户文件列表
func QueryFileHandler(w http.ResponseWriter, r *http.Request) {
	// 获取用户名
	username := r.Form.Get("username")
	userFiles, err := db.QueryUserFileDB(username, 10)

	var resp util.RespMsg
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp = util.RespMsg{Msg: "Failed get file meta , err : " + err.Error()}
	} else {
		w.WriteHeader(http.StatusOK)
		resp = util.RespMsg{
			Msg:  "通过token获取文件元信息",
			Data: userFiles,
		}
	}

	util.Logerr(w.Write(resp.JSONBytes()))
}

// 文件删除的接口
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}
	fileSha1 := r.Form.Get("filehash")

	fMeta, err := db.GetFileDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fMeta.Location)
	//meta.RemoveFileMeta(fileSha1)

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

	curFileMeta, err := db.GetFileDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curFileMeta.FileName = newFileName
	db.UpdateFileDB(curFileMeta.FileName, curFileMeta.FileSha1)

	resp := util.RespMsg{
		Msg:  "通过sha1更新文件元信息",
		Data: curFileMeta,
	}

	w.WriteHeader(http.StatusOK)

	util.Logerr(w.Write(resp.JSONBytes()))
}

// 下载文件的接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	fsha1 := r.Form.Get("filehash")
	fm, err := db.GetFileDB(fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
