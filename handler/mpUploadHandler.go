package handler

import (
	"filestore/cache/myredis"
	"fmt"
	"net/http"
)

// 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求参数
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: "+err.Error())
		return
	}

	//2. 获得一个Redis连接
	rConn := myredis.RedisPool().Get()
	defer rConn.Close()

	//3. 生成分块上传的初始化信息
	//upload := MultipartUploadInfo{
	//	FileHash:   "",
	//	FileSize:   0,
	//	UploadID:   "",
	//	ChunkSize:  0,
	//	ChunkCount: 0,
	//}

	//4. 将初始化信息写入到redis缓存
	//5. 将响应初始化数据返回到客户端
}
