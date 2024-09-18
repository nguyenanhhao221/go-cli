package main

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"testing"
	"testing/iotest"
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

func TestCsv2Float(t *testing.T) {
	csvData := `IP Address,Timestamp,Response Time,Bytes
192.168.0.199,1520698621,236,3475
192.168.0.88,1520698776,220,3200
192.168.0.199,1520699033,226,3200
192.168.0.100,1520699142,218,3475
192.168.0.199,1520699379,238,3822
`
	testCases := []struct {
		name           string
		col            int
		expectedResult []float64
		expErr         error
		r              io.Reader
	}{
		{name: "Column2", col: 2, expectedResult: []float64{1520698621, 1520698776, 1520699033, 1520699142, 1520699379}, expErr: nil, r: bytes.NewBufferString(csvData)},
		{name: "Column3", col: 3, expectedResult: []float64{236, 220, 226, 218, 238}, expErr: nil, r: bytes.NewBufferString(csvData)},
		{name: "FailRead", col: 1, expErr: iotest.ErrTimeout, r: iotest.TimeoutReader(bytes.NewReader([]byte{0}))},
		{name: "Fail Invalid Column", col: 5, expErr: ErrInvalidColumn, r: bytes.NewBufferString(csvData)},
		{name: "Fail Not Number", col: 1, expErr: ErrNotNumber, r: bytes.NewBufferString(csvData)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := csv2float(tc.r, tc.col)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expect error, got nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q got %q", tc.expErr, err)
				}
			}
			if !slices.Equal(result, tc.expectedResult) {
				t.Errorf("Expected %#v, got %#v", tc.expectedResult, result)
			}
		})

	}
}
