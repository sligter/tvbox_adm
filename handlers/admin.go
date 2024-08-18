package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"tvbox/models"

	"github.com/gorilla/mux"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	links, err := models.GetAllLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hasActiveLink := false
	for _, link := range links {
		if link.State == 1 {
			hasActiveLink = true
			break
		}
	}

	data := struct {
		Links         []models.Link
		HasActiveLink bool
	}{
		Links:         links,
		HasActiveLink: hasActiveLink,
	}

	// 使用预先解析的模板
	err = templates.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	name := r.FormValue("name")
	state, _ := strconv.Atoi(r.FormValue("status"))

	err := models.AddLink(url, name, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteLink(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ChangeStateHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = models.ChangeState(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	name := r.FormValue("name")

	err = models.UpdateLink(id, url, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func TVBoxHandler(w http.ResponseWriter, r *http.Request) {
	url, err := models.GetActiveLink()
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到活跃链接，重定向到管理界面
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		// 其他错误则返回内部服务器错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果找到活跃链接，则重定向到该链接
	http.Redirect(w, r, url, http.StatusSeeOther)
}
