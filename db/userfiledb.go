package db

import (
	"filestore-server/db/mydb"
	"fmt"
)

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
