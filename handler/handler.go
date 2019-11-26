package handler

import (
	"filestore-server/meta"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// 下载文件的接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	fsha1 := r.Form.Get("filehash")
	fm, err := meta.GetFileMetaDB(fsha1)
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
