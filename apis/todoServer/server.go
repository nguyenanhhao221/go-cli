package main

import (
	"log"
	"net/http"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("GET /", rootHandler)
	return m
}

func replyPlainText(w http.ResponseWriter, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(content)); err != nil {
		log.Printf("error writing response: %v", err) // Log the error if something goes wrong
	}
}
