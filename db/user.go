package db

import (
	"database/sql"
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

// 判断密码是否正确
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT user_pwd FROM fileserver_user WHERE user_name=?",
	)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	var encpwddb string
	err = stmt.QueryRow(username).Scan(&encpwddb)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Failed to select, err : " + err.Error())
		}
		return false
	}

	if encpwddb == encpwd {
		return true
	} else {
		return false
	}
}

type TableUser struct {
	Uid        sql.NullInt64
	Username   sql.NullString
	Email      sql.NullString
	Phone      sql.NullString
	SignUpAt   string
	LastAction string
}

// 获取用户信息
func GetUserInfo(username string) (*TableUser, error) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT uid, user_name , email, phone, signup_at, last_active FROM fileserver_user WHERE user_name=?",
	)
	if err != nil {
		fmt.Println("Failed to select, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	var tuser TableUser
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
