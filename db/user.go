package db

import (
	"filestore-server/db/mydb"
	"fmt"
)

// 通过用户名以及密码完成user表的注册操作
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT IGNORE INTO fileserver_user(`user_name`, `user_pwd`) VALUES (?, ?)",
	)
	if err != nil {
		fmt.Println("Filed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Filed to insert, err:" + err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("User has been uploaded before\n")
			return false
		}
		return true
	}
	return false
}
