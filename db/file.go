package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// 用于更新文件元信息的数据库操作封装类
func OnFileUploadFinish(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT IGNORE INTO fileserver_file (`file_sha1`, `file_name`, `file_size`, " +
			" `file_addr`,  `status` ) VALUES(?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("File with hash:%s has been uploaded before", filehash)
		}

		return true
	}
	return false
}
