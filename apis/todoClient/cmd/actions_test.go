package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
		isActive    bool
	}{
		{
			name:     "Results",
			expError: nil,
			expOut:   "-  1  Task 1\n-  2  Task 2\n",
			resp:     testResp["resultsMany"],
		},
		{
			name:     "NoResults",
			expError: ErrNotFound,
			resp:     testResp["noResults"],
		},
		{
			name:        "InvalidURL",
			expError:    ErrConnection,
			resp:        testResp["noResults"],
			closeServer: true,
		},
		{
			name:     "ResultOnlyActive",
			expError: nil,
			resp:     testResp["resultsManyComplete"],
			expOut:   "-  2  Task 2\n",
			isActive: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})

			defer cleanup()

			if tc.closeServer {
				cleanup()
			}

			var out bytes.Buffer

			err := listAction(&out, url, tc.isActive)
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, no error", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expect error %q, got %q", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error got %q", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expect output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		id       string
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:     "ResultOne",
			id:       "1",
			expError: nil,
			expOut: `Task:         Task 1
Created at:   Oct/28 @08:23
Completed:    No
`,
			resp: testResp["resultsOne"],
		},
		{
			name:     "NoResults",
			id:       "1",
			expError: ErrNotFound,
			expOut:   ``,
			resp:     testResp["notFound"],
		},
		{
			name:     "InvalidID",
			id:       "foo",
			expError: ErrNotNumber,
			expOut:   ``,
			resp:     testResp["noResults"],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})

			defer cleanup()

			if tc.closeServer {
				cleanup()
			}
			var out bytes.Buffer
			err := viewAction(&out, url, tc.id)
			if tc.expError != nil {
				if err == nil {
					t.Error("Exp to have error got nil\n")
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expect error %v, got %v \n", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expect %s got %s\n", tc.expOut, out.String())
			}
		})
	}
}

// TestAddAction tests the addAction function using a mock server.
// It verifies that addAction sends the correct request (URL, method, body, and headers)
// without actually calling the real API. This ensures that our logic is correct
// before interacting with the actual Task Todo API.
func TestAddAction(t *testing.T) {
	expURLPath := "/todo"
	expMethod := http.MethodPost
	expBody := "{\"task\":\"Task 1\"}\n"
	expContentType := "application/json"
	expOut := "Added task \"Task 1\" to the list.\n"
	args := []string{"Task", "1"}

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURLPath {
			t.Errorf("Expected Path %q, got %q", expURLPath, r.URL.Path)
		}
		if r.Method != expMethod {
			t.Errorf("Expect POST Method, got %q", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != expBody {
			t.Errorf("Expected body %q, got %q", expBody, string(body))
		}

		if r.Header.Get("Content-Type") != expContentType {
			t.Errorf("Expected header to have Content-Type:  %q, got %q", r.Header.Get("Content-Type"), expContentType)
		}
		w.WriteHeader(testResp["created"].Status)
		fmt.Fprintln(w, testResp["created"].Body)
	})

	defer cleanup()

	var out bytes.Buffer

	err := addAction(&out, url, args)
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != expOut {
		t.Errorf("Expect %q, got %q", expOut, out.String())
	}
}

func TestCompleteAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodPatch
	expQuery := "complete"
	expOut := "Item number 1 marked as completed.\n"
	arg := "1"

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURLPath {
			t.Errorf("Expected Path %q, got %q\n", expURLPath, r.URL.Path)
		}
		if r.Method != expMethod {
			t.Errorf("Expect %q Method, got %q\n", expMethod, r.Method)
		}
		if r.URL.Query().Get("complete") != "" {
			t.Errorf("Expect query to be %q, got %q\n", expQuery, r.URL.Query().Get("complete"))
		}

		w.WriteHeader(testResp["noContent"].Status)
		fmt.Fprintln(w, testResp["noContent"].Body)
	})
	defer cleanup()

	var out bytes.Buffer
	err := completeAction(&out, url, arg)
	if err != nil {
		t.Fatal(err)
	}

	if expOut != out.String() {
		t.Errorf("Expect %q, got %q", expOut, out.String())
	}
}

func TestDeleteAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodDelete
	expOut := "Item number 1 deleted\n"
	arg := "1"

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURLPath {
			t.Errorf("Expected Path %q, got %q\n", expURLPath, r.URL.Path)
		}
		if r.Method != expMethod {
			t.Errorf("Expect %q Method, got %q\n", expMethod, r.Method)
		}

		w.WriteHeader(testResp["noContent"].Status)
		fmt.Fprintln(w, testResp["noContent"].Body)
	})
	defer cleanup()

	var out bytes.Buffer
	err := deleteAction(&out, url, arg)
	if err != nil {
		t.Fatal(err)
	}

	if expOut != out.String() {
		t.Errorf("Expect %q, got %q", expOut, out.String())
	}
}
