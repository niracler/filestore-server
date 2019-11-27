package meta

// 文件元信息结构
type FileMeta struct {
	FileSha1 string `json:"fileSha1"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	Location string `json:"location"`
	UploadAt string `json:"uploadAt"`
}
