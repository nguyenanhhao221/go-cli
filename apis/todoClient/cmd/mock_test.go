package cmd

import (
	"net/http"
	"net/http/httptest"
)

// testResp simulates test responses from the API
var testResp = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body: `{
  "results": [
    {
      "Task": "Task 1",
      "Done": false,
      "CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
      "CompletedAt": "0001-01-01T00:00:00Z"
    },
    {
      "Task": "Task 2",
      "Done": false,
      "CreatedAt": "2019-10-28T08:23:38.323447798-04:00",
      "CompletedAt": "0001-01-01T00:00:00Z"
    }
  ],
  "date": 1572265440,
  "totalResults": 2
}`,
	},
	"resultsOne": {
		Status: http.StatusOK,
		Body: `{
  "results": [
    {
      "Task": "Task 1",
      "Done": false,
      "CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
      "CompletedAt": "0001-01-01T00:00:00Z"
    }
  ],
  "date": 1572265440,
  "totalResults": 1
}`,
	},

	"noResults": {
		Status: http.StatusOK,
		Body: `{
  "results": [],
  "date": 1572265440,
  "totalResults": 0
}`,
	},

	"resultsManyComplete": {
		Status: http.StatusOK,
		Body: `{
  "results": [
    {
      "Task": "Task 1",
      "Done": true,
      "CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
      "CompletedAt": "0001-01-01T00:00:00Z"
    },
    {
      "Task": "Task 2",
      "Done": false,
      "CreatedAt": "2019-10-28T08:23:38.323447798-04:00",
      "CompletedAt": "0001-01-01T00:00:00Z"
    }
  ],
  "date": 1572265440,
  "totalResults": 2
}`,
	},
	"root": {
		Status: http.StatusOK,
		Body:   "There's an API here",
	},

	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},

	"created": {
		Status: http.StatusCreated,
		Body:   "",
	},
	"noContent": {
		Status: http.StatusNoContent,
		Body:   "",
	},
}

// mockServer creates a mock server to simulate the todo API
func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() {
		ts.Close()
	}
}
