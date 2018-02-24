package main

import (
	"net/http"
	"log"
	"html/template"
	"io/ioutil"
	"path"
	"runtime/debug"
	"LBSWeb/session"
	_ "LBSWeb/session/memory"
	"os"
)

var userMgr *UserManager
var globalSessions *session.Manager

const (
	TEMPLATE_DIR = "./views"
	STATIC_DIR = "./static"
)

var templates map[string]*template.Template

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok :=  recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// RenderHtml(w, "error", e)
				// logging
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func RenderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	 key := TEMPLATE_DIR + "/" + tmpl + ".html"
	 err := templates[key].Execute(w, locals)
	 check(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if !globalSessions.SessionExist(r) {
			log.Println("session is not exist!!!")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		//跳转到主页服务
		RenderHtml(w, "index", nil)
	}
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func staticDirHandler(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := STATIC_DIR + r.URL.Path[len(prefix)-1:]
		log.Println("static request:", file)
		if exists := IsExist(file); !exists {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, file)
	})
}

func init() {
	userMgr = NewUserManager()
	// for test init users
	userMgr.Users["root"] = "root"

	globalSessions, _ = session.NewManager("memory", "xmtsessionid", 3600)
	go globalSessions.GC()

	templates = make(map[string]*template.Template)
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		check(err)
		return
	}

	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
	}
}

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/")
	mux.HandleFunc("/", safeHandler(homeHandler))
	mux.HandleFunc("/login", safeHandler(userMgr.LoginHandler))
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
