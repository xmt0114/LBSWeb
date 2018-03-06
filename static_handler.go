package main

import (
	"net/http"
	"LBSWeb/logger"
)

const STATIC_DIR = "E:/code/go/src/LBSWeb/static"

func staticDirHandler(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := STATIC_DIR + r.URL.Path[len(prefix)-1:]
		logger.Logger.Debugf("static request:", file)
		if exists := IsExist(file); !exists {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, file)
	})
}
