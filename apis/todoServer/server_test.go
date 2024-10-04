package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"todo"
)

func TestServerGet(t *testing.T) {
	testCases := []struct {
		name       string
		route      string
		expStatus  int
		expItems   int
		expReponse string
	}{
		{name: "GetRoute", route: "/", expStatus: http.StatusOK, expReponse: "hello there you hit the api"},
		{name: "GetAll", route: "/todo", expStatus: http.StatusOK, expItems: 2, expReponse: "Task number 1."},
		{name: "GetSingle", route: "/todo/2", expStatus: http.StatusOK, expItems: 1, expReponse: "Task number 2."},
		{name: "NotFoundRoute", route: "/invalid/todo", expStatus: http.StatusNotFound, expReponse: "404 page not found\n"},
	}

	serverUrl, cleanup := setupTestServer(t)
	defer cleanup()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.Get(serverUrl + tc.route)
			if err != nil {
				t.Error(err)
			}
			var (
				body []byte
				resp struct {
					Results      todo.List `json:"results"`
					Date         int64     `json:"date"`
					TotalResults int       `json:"totalResults"`
				}
			)
			if r.StatusCode != tc.expStatus {
				t.Errorf("Expect status: %d, got %d\n", tc.expStatus, r.StatusCode)
			}
			defer r.Body.Close()

			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}
				if string(body) != tc.expReponse {
					t.Errorf("Expect response %q, got :%q\n", tc.expReponse, string(body))
				}
			case strings.Contains(r.Header.Get("Content-Type"), "application/json"):
				if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Error(err)
				}

				if resp.TotalResults != tc.expItems {
					t.Errorf("Expect totalResults: %d, got %d\n", tc.expItems, resp.TotalResults)
				}

				if resp.Results.Items[0].Task != tc.expReponse {
					t.Errorf("Expect task name: %q, got %q\n", tc.expReponse, resp.Results.Items[0].Task)
				}
			}

		})
	}
}

func TestAddTodo(t *testing.T) {
	serverUrl, cleanup := setupTestServer(t)
	defer cleanup()

	task := struct {
		Task string `json:"task"`
	}{
		Task: "foo",
	}
	jTask, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.Post(serverUrl+"/todo", "application/json", bytes.NewBuffer(jTask))
	if err != nil {
		t.Fatal(err)
	}

	defer r.Body.Close()
	if r.StatusCode != http.StatusCreated {
		t.Errorf("Expect status %d, got %d", http.StatusCreated, r.StatusCode)
	}
}

func setupTestServer(t *testing.T) (string, func()) {
	t.Helper()

	var tempFile, err = os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(newMux(tempFile.Name()))
	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d.", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}
	return ts.URL, func() {
		ts.Close()
	}
}
