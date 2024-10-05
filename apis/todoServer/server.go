package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"todo"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	mu := &sync.Mutex{}
	list := &todo.List{}

	m.HandleFunc("GET /", rootHandler)
	m.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
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

		mu.Lock()
		defer mu.Unlock()
		getSingleTodoRouter(w, r, list, id)
	})

	//POST
	m.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		addTodoRouter(w, r, list, todoFile)
	})

	// DELETE
	m.HandleFunc("DELETE /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := validateID(r.PathValue("id"), list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		mu.Lock()
		defer mu.Unlock()
		deleteTodoHandler(w, r, list, id, todoFile)
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
