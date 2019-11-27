package db

import (
	"filestore-server/db/mydb"
	"filestore-server/meta"
	"fmt"
)

// 创建用户文件
func CreateUserFileDB(username string, filehash string, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT IGNORE INTO fileserver_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		fmt.Println("准备语句有问题:", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize)
	if err != nil {
		fmt.Println("执行语句有问题:" + err.Error())
		return false
	}

	return true
}

// 批量获取用户文件信息
func QueryUserFileDB(username string, limit int) ([]meta.UserFileMeta, error) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT file_sha1, file_size, file_name, upload_at, last_update FROM fileserver_user_file WHERE user_name=? LIMIT ?",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []meta.UserFileMeta
	for rows.Next() {
		ufile := meta.UserFileMeta{}
		ufile.Username = username
		err = rows.Scan(&ufile.FileHash, &ufile.FileSize, &ufile.FileName, &ufile.SignUpAt, &ufile.LastAction)
		if err != nil {
			return nil, err
		}
		userFiles = append(userFiles, ufile)
	}

	return userFiles, nil
}
