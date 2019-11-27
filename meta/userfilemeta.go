package meta

// 用户文件元信息
type UserFileMeta struct {
	Username   string `json:"username"`
	FileHash   string `json:"fileHash"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	SignUpAt   string `json:"signUpAt"`
	LastAction string `json:"lastAction"`
}
