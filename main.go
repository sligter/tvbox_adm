package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"tvbox/db"
	"tvbox/handlers"

	"github.com/gorilla/mux"
)

//go:embed static/*
var staticContent embed.FS

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "0.0.0.0"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	// 定义命令行参数
	port := flag.Int("port", 2345, "Port to run the server on")
	domain := flag.String("domain", "", "Domain to bind to (default: local IP)")
	flag.Parse()

	// 初始化数据库
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	handlers.InitHandlers(staticContent)
	// 路由设置
	r.HandleFunc("/", handlers.TVBoxHandler)
	r.HandleFunc("/admin", handlers.AdminHandler)
	r.HandleFunc("/admin/login", handlers.LoginHandler)
	r.HandleFunc("/admin/logout", handlers.LogoutHandler)
	r.HandleFunc("/admin/add-link", handlers.AddLinkHandler)
	r.HandleFunc("/admin/edit-link", handlers.EditLinkHandler)
	r.HandleFunc("/admin/delete-link/{id}", handlers.DeleteLinkHandler)
	r.HandleFunc("/admin/change-state/{id}", handlers.ChangeStateHandler)
	r.HandleFunc("/admin/update-password", handlers.UpdatePasswordHandler)

	// 添加静态文件服务
	fileServer := http.FileServer(http.FS(staticContent))
	r.PathPrefix("/static/").Handler(http.StripPrefix("", fileServer))

	// 构建服务器地址
	localIP := getLocalIP()
	listenAddr := fmt.Sprintf(":%d", *port)

	var publicURL string
	if *domain != "" {
		publicURL = fmt.Sprintf("http://%s:%d", *domain, *port)
	} else {
		publicURL = fmt.Sprintf("http://%s:%d", localIP, *port)
	}

	log.Printf("Server starting, public URL: %s\n", publicURL)
	log.Printf("Server listening on %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
