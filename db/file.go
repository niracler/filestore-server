package db

import (
	"database/sql"
	"filestore-server/db/mydb"
	"fmt"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	FileCre  string
}

// 用于更新文件元信息的数据库操作封装类
func CreateFileDB(filehash string, filename string, filesize int64, fileaddr string) bool {
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
			fmt.Printf("File with hash:%s has been uploaded before\n", filehash)
		}

		return true
	}
	return false
}

// 用于更新文件元信息的数据库操作封装类
func UpdateFileDB(filehash string, filename string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"UPDATE fileserver_file SET file_name=? WHERE file_sha1=?",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filename, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before\n", filehash)
		}

		return true
	}
	return false
}

// 从MySQL获取文件元信息
func GetFileDB(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT file_sha1, file_addr, file_name, file_size, created FROM fileserver_file WHERE file_sha1=? AND status=1",
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize, &tfile.FileCre)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tfile, nil
}
