package main

import (
	"sync"
	"net/http"
	"log"
)

type UserManager struct {
	lock  sync.RWMutex
	Users map[string]string
}

func NewUserManager() *UserManager {
	return &UserManager{Users:make(map[string]string)}
}

func init()  {
	// for test init users
	//TODO: load user db
	userMgr.Users["root"] = "root"
}

func (mgr *UserManager) CheckUser(name string, passwd string) bool {
	if name == "" || passwd == "" {
		return false
	}

	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	if mgr.Users[name] == passwd {
		return true
	}
	return  false
}

func (mgr *UserManager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("user login by ", r.Method)
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "GET" {
		log.Println("user login access login page")
		info := make(map[string]interface{})
		info["info"] = "请您登录"
		RenderHtml(w, "login", nil)
		return
	} else if r.Method == "POST" {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if !mgr.CheckUser(username, password) {
			log.Println("user login but error")
			info := make(map[string]interface{})
			info["info"] = "用户名或密码错误"
			RenderHtml(w, "login", info)
			return
		}

		sess.Set("username", username)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}