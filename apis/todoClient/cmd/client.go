package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("Connection Error")
	ErrNotFound        = errors.New("Not Found")
	ErrInvalidResponse = errors.New("Invalid Server Response")
	ErrInvalid         = errors.New("Invalid Data")
	ErrNotNumber       = errors.New("Not a number")
)

const timeFormat = "Jan/02 @15:04"

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type todoResponse struct {
	Results      []item `json:"results"`
	Date         int64  `json:"date"`
	TotalResults int    `json:"totalResults"`
}

func newClient() *http.Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client
}

func getItems(endpoint string) ([]item, error) {
	client := newClient()
	res, err := client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Cannot read body %w", err)
		}

		err = ErrInvalidResponse
		if res.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	var resp todoResponse

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("%w: fail to decode json response", err)
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w: No results found", ErrNotFound)
	}
	return resp.Results, nil
}

func getOne(endpoint string) (item, error) {
	items, err := getItems(endpoint)
	if err != nil {
		return item{}, err

	}
	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}
	return items[0], nil
}
