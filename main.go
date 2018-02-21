package main

import (
	"net/http"
	"log"
	"html/template"
	"io/ioutil"
	"path"
	"runtime/debug"
)

const (
	TEMPLATE_DIR = "E:/code/go/src/LBSWeb/views"
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
				// renderHtml(w, "error", e)
				// logging
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	 key := TEMPLATE_DIR + "/" + tmpl + ".html"
	 err := templates[key].Execute(w, locals)
	 check(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderHtml(w, "index", nil)
}

func init() {
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
	http.HandleFunc("/", safeHandler(indexHandler))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
