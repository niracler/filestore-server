package meta

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

// 文件元信息初始化
func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增、更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 通过sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 通过sha1删除文件
func RemoveFileMeta(fileSha1 string) {

}
