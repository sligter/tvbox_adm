package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./tvbox.db")
	if err != nil {
		return err
	}

	// 创建 tvbox 表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tvbox (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT,
        state INTEGER,
        name TEXT
    )`)
	if err != nil {
		return err
	}

	// 创建 admin 表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS admin (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        password TEXT,
        auth_token TEXT
    )`)
	if err != nil {
		return err
	}

	// 检查是否已存在管理员账户
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM admin").Scan(&count)
	if err != nil {
		return err
	}

	// 如果不存在管理员账户，创建初始账户
	if count == 0 {
		password := hashPassword("admin")
		_, err = DB.Exec("INSERT INTO admin (username, password) VALUES (?, ?)", "admin", password)
		if err != nil {
			return err
		}
	}

	return nil
}

func hashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
