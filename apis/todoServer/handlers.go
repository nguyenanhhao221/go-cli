package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	content := "hello there you hit the api"
	replyPlainText(w, http.StatusOK, content)
}
