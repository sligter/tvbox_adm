package models

import (
	"tvbox/db"
)

func ValidateAdmin(username, password string) bool {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM admin WHERE username = ? AND password = ?", username, password).Scan(&count)
	return err == nil && count > 0
}

func UpdateAuthToken(username, token string) error {
	_, err := db.DB.Exec("UPDATE admin SET auth_token = ? WHERE username = ?", token, username)
	return err
}

func ValidateAuthToken(token string) bool {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM admin WHERE auth_token = ?", token).Scan(&count)
	return err == nil && count > 0
}

func UpdateAdminPassword(oldPassword, newPassword string) error {
	_, err := db.DB.Exec("UPDATE admin SET password = ? WHERE password = ?", newPassword, oldPassword)
	return err
}
