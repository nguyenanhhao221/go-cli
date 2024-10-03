package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerGet(t *testing.T) {
	testCases := []struct {
		name       string
		route      string
		expStatus  int
		expReponse string
	}{
		{name: "GetRoute", route: "/", expStatus: http.StatusOK, expReponse: "hello there you hit the api"},
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
			}
		})
	}

}

func setupTestServer(t *testing.T) (string, func()) {
	t.Helper()

	ts := httptest.NewServer(newMux(""))
	return ts.URL, func() {
		ts.Close()
	}
}
