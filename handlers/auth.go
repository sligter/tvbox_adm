package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"tvbox/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if models.ValidateAdmin(username, hashPassword(password)) {
			token := generateToken()
			err := models.UpdateAuthToken(username, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "auth_token",
				Value:   token,
				Expires: time.Now().Add(7 * 24 * time.Hour),
			})

			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		errorMessage = "Username or password error"
	}

	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}

	err := templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "auth_token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirmNewPassword := r.FormValue("confirm_new_password")

	if newPassword != confirmNewPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	err := models.UpdateAdminPassword(hashPassword(oldPassword), hashPassword(newPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return false
	}

	return models.ValidateAuthToken(cookie.Value)
}

func hashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateToken() string {
	return hashPassword(time.Now().String())
}
