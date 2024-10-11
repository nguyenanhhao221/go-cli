package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo"
)

var (
	ErrInvalidData error = errors.New("invalid data")
	ErrNotFound    error = errors.New("not found")
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	content := "hello there you hit the api"
	replyPlainText(w, r, http.StatusOK, content)
}

func getTodoRouter(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {
	if err := list.Get(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
	}

	resp := &todoResponse{
		Results: list.Items,
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func getSingleTodoRouter(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {
	resp := &todoResponse{
		Results: list.Items[id-1 : id],
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func addTodoRouter(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {
	// Add todo
	item := struct {
		Task string `json:"task"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		replyError(w, r, http.StatusBadRequest, err.Error())
	}

	list.Add(item.Task)

	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	replyPlainText(w, r, http.StatusCreated, "")
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	if err := list.Delete(id); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	replyPlainText(w, r, http.StatusNoContent, "")
}

func pacthHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	if err := list.Complete(id); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	replyPlainText(w, r, http.StatusNoContent, "")
}

func validateID(idString string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID :%s", ErrInvalidData, err)
	}

	if id < 1 {
		return 0, fmt.Errorf("%w: Invalid ID: less than 1", ErrInvalidData)
	}

	if id > len(list.Items) {
		return 0, fmt.Errorf("%w Id not found: %d", ErrNotFound, id)
	}

	return id, nil
}
