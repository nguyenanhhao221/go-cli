package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name      string
		col       int
		filenames []string
		op        string
		expected  string
		expErr    error
	}{
		{name: "RunAvg1File", op: "avg", col: 3, filenames: []string{"testdata/example.csv"}, expected: "227.6\n", expErr: nil},
		{name: "RunAvgMultiFiles", op: "avg", col: 3, filenames: []string{"testdata/example.csv", "testdata/example2.csv"}, expected: "233.84\n", expErr: nil},
		{name: "RunFailReadFile", op: "avg", col: 3, filenames: []string{"testdata/invalid.csv"}, expErr: os.ErrNotExist},
		{name: "RunFailColumn", op: "avg", col: 0, filenames: []string{"testdata/example.csv"}, expErr: ErrInvalidColumn},
		{name: "RunFailNoFiles", op: "avg", col: 0, filenames: []string{}, expErr: ErrNoFiles},
		{name: "RunFailOperation", op: "foo", col: 3, filenames: []string{"testdata/example.csv"}, expErr: ErrInvalidOperation},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer
			err := run(tc.filenames, tc.op, tc.col, &res)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expect error, got nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Expect no error, got %q", err)
			}

			out := res.String()
			if tc.expected != out {
				t.Errorf("Expected %q, got %q", tc.expected, &res)
			}
		})
	}
}
