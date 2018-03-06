package main

import (
	"net/http"
	"html/template"
	"LBSWeb/session"
	_ "LBSWeb/session/memory"
	"LBSWeb/logger"
)

var userMgr *UserManager
var globalSessions *session.Manager
var templates map[string]*template.Template

func init() {
	logger.Logger.Info("LBSWeb init")
	userMgr = NewUserManager()

	globalSessions, _ = session.NewManager("memory", "xmtsessionid", 3600)
	go globalSessions.GC()

	LoadTmplates()
}

func main() {
	logger.Logger.Info("LBSWeb start")
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/")
	mux.HandleFunc("/", SafeHandler(HomeHandler))
	mux.HandleFunc("/login", SafeHandler(userMgr.LoginHandler))

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		logger.Logger.Critical("ListenAndServe: ", err.Error())
	}
}
