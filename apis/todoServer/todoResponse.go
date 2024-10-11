package main

import (
	"encoding/json"
	"time"
	"todo"
)

type todoResponse struct {
	Results []todo.Item `json:"results"`
}

type item struct {
	CompletedAt time.Time
	CreatedAt   time.Time
	Task        string
	Done        bool
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      []todo.Item `json:"results"`
		Date         int64       `json:"date"`
		TotalResults int         `json:"totalResults"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}

	return json.Marshal(resp)
}
