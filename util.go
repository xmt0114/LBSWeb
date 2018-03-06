package main

import (
	"net/http"
	"LBSWeb/logger"
	"runtime/debug"
	"html/template"
	"os"
	"fmt"
	"io/ioutil"
	"path"
)

const TEMPLATE_DIR = "E:/code/go/src/LBSWeb/views"

func LoadTmplates()  {
	templates = make(map[string]*template.Template)
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		Check(err)
		return
	}

	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		logger.Logger.Debugf("Loading template: %v", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
	}
}

func RenderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	key := TEMPLATE_DIR + "/" + tmpl + ".html"
	if v, ok := templates[key]; ok {
		err := v.Execute(w, locals)
		Check(err)
	} else {
		logger.Logger.Warnf("Not Found matched template for %s", tmpl)
		Check(fmt.Errorf("Not Found matched template for %s", tmpl))
	}
}

func CheckSessionAndRedirectLogPage(w http.ResponseWriter, r *http.Request) bool {
	if !globalSessions.SessionExist(r) {
		logger.Logger.Trace("session is not exist!!!")
		http.Redirect(w, r, "/login", http.StatusFound)
		return false
	}
	return true
}

func SafeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok :=  recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// RenderHtml(w, "error", e)
				// logging
				logger.Logger.Warnf("WARN: panic in %v - %v", fn, e)
				logger.Logger.Warn(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}