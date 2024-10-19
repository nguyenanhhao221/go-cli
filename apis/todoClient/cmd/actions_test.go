package cmd

import (
	"bytes"
	"errors"
	"fmt"
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

			err := listAction(&out, url)
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
