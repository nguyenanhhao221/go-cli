package main

import (
	"testing"
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}
	testCases := []struct {
		name     string
		op       statsFunc
		expected []float64
	}{
		{"Sum", sum, []float64{300, 85.927, -30, 436}},
		{"Average", avg, []float64{37.5, 6.609769230769231, -15, 72.666666666666666}},
	}

	for _, tc := range testCases {
		for i, v := range data {
			t.Run(tc.name, func(t *testing.T) {
				result := tc.op(v)
				// Comparing the float is quite tricky, temporary compare like this
				if result != tc.expected[i] {
					t.Errorf("Expected %g got %g", tc.expected[i], result)
				}
			})
		}
	}
}
