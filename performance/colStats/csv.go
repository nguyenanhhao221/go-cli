package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strconv"
)

func sum(data []float64) float64 {
	sum := 0.0

	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func min(data []float64) float64 {
	return slices.Min(data)
}

func max(data []float64) float64 {
	return slices.Max(data)
}

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

func csv2float(r io.Reader, column int) ([]float64, error) {
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true
	// Adjusting for 0 based index
	column--

	var data []float64

	for i := 0; ; i++ {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Cannot read data from file: %w", err)
		}

		if i == 0 {
			continue
		}

		if len(row) <= column {
			// File does not have that many column
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		// Try to convert data to float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}
	return data, nil
}
