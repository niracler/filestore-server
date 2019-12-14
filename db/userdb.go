package db

import (
	"database/sql"
	"filestore/db/mydb"
	"filestore/meta"
	"fmt"
)

// 通过用户名以及密码完成user表的注册操作
func CreateUserDB(username string, password string) bool {
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

// 判断密码是否正确
func IsValidUserDB(username string, encpwd string) (int64, bool) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT uid, user_pwd FROM fileserver_user WHERE user_name=?",
	)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return 0, false
	}
	defer stmt.Close()

	var encpwddb string
	var uid int64
	err = stmt.QueryRow(username).Scan(&uid, &encpwddb)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Failed to select, err : " + err.Error())
		}
		return 0, false
	}

	if encpwddb == encpwd {
		return uid, true
	} else {
		return 0, false
	}
}

// 获取用户信息
func GetUserInfoDB(username string) (*meta.UserMeta, error) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT uid, user_name , email, phone, signup_at, last_active FROM fileserver_user WHERE user_name=?",
	)
	if err != nil {
		fmt.Println("Failed to select, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	var tuser meta.UserMeta
	err = stmt.QueryRow(username).Scan(&tuser.Uid, &tuser.Username, &tuser.Email, &tuser.Phone, &tuser.SignUpAt, &tuser.LastAction)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Failed to select, err : " + err.Error())
		}
		return nil, err
	}
	return &tuser, nil
}
