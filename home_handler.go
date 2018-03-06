package main

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if !CheckSessionAndRedirectLogPage(w, r) {
			return
		}

		//跳转到主页服务
		RenderHtml(w, "index", nil)
	}
}
