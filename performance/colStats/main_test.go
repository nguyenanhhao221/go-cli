package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
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
		{name: "RunMin1File", op: "min", col: 3, filenames: []string{"testdata/example.csv"}, expected: "218\n", expErr: nil},
		{name: "RunMax1File", op: "max", col: 3, filenames: []string{"testdata/example.csv"}, expected: "238\n", expErr: nil},
		{name: "RunAvgMultiFiles", op: "avg", col: 3, filenames: []string{"testdata/example.csv", "testdata/example2.csv"}, expected: "233.84\n", expErr: nil},
		{name: "RunSumMultiFiles", op: "sum", col: 3, filenames: []string{"testdata/example.csv", "testdata/example2.csv"}, expected: "5846\n", expErr: nil},
		{name: "RunMinMultiFiles", op: "min", col: 3, filenames: []string{"testdata/example.csv", "testdata/example2.csv"}, expected: "218\n", expErr: nil},
		{name: "RunMaxMultiFiles", op: "max", col: 3, filenames: []string{"testdata/example.csv", "testdata/example2.csv"}, expected: "238\n", expErr: nil},
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

func BenchmarkRun(b *testing.B) {
	filesname, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	// In Go benchmark tests, the b.ResetTimer() function is important because it resets the internal timer that Go uses to measure the time taken for a benchmark to run. Here’s why this is crucial:
	//
	// 	1.	Excludes Setup Time: When you run a benchmark, there may be some setup or initialization code that runs before the actual code you want to measure. Without ResetTimer(), the setup time would be included in the benchmark’s result, which would give an inaccurate measure of the code’s actual performance.
	// 	2.	Accurate Timing of Core Logic: By calling b.ResetTimer() just before the main benchmark loop, you ensure that the timer only captures the time taken for the core logic you are interested in benchmarking, excluding any setup, file loading, or pre-calculation steps.
	//
	// For example, if you’re reading files or preparing data before running the benchmarked function, ResetTimer ensures those steps aren’t included in the performance metrics.
	//
	// In summary, ResetTimer() ensures you’re measuring the performance of only the part of your code you want to evaluate, leading to more accurate benchmarking results.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filesname, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunMin(b *testing.B) {
	filesname, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filesname, "min", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
