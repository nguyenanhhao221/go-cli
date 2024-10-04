package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"todo"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	list := &todo.List{}

	m.HandleFunc("GET /", rootHandler)
	m.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		getTodoRouter(w, r, list, todoFile)
	})
	m.HandleFunc("GET /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := validateID(r.PathValue("id"), list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		getSingleTodoRouter(w, r, list, id)
	})

	m.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		addTodoRouter(w, r, list, todoFile)
	})
	return m
}

func replyPlainText(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(content)); err != nil {
		log.Printf("error writing response: %v", err)
		replyError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	body, err := json.Marshal(resp)
	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func replyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, message)
	http.Error(w, message, status)
}
