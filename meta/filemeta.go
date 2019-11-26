package meta

import "filestore-server/db"

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// 新增文件元信息到 mysql 中
func CreateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinish(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 更新文件元信息到 mysql 中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUpdateFinish(fmeta.FileSha1, fmeta.FileName)
}

// 通过sha1从数据库获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}

	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.FileCre,
	}

	return fmeta, nil
}

// 通过sha1删除文件
func RemoveFileMeta(fileSha1 string) {

}
