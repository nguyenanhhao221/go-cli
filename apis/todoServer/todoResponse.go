package main

import (
	"encoding/json"
	"time"
	"todo"
)

type todoResponse struct {
	Results todo.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"totalResults"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results.Items),
	}

	return json.Marshal(resp)
}