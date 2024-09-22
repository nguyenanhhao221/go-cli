package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	// Validate operation input from user
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	wg := sync.WaitGroup{}

	errCh := make(chan error)
	resCh := make(chan []float64)
	doneCh := make(chan struct{})

	for _, fname := range filenames {
		wg.Add(1)
		// for each file, we use an anonymous go routine,
		// We pass the fname again to this function to avoid common bug for Go before 1.22 due to the way for loop work and close.
		// this is not necessary for go >1.22
		// Before Go 1.22, if we don't pass fname as param to this function, each time Go went through the loop, it create and reuse the same variable
		// In this case, when we go through filenames, it create fname variable, then next time, it reuse that same variable
		// Because of this, if use go routine like we did, we will always get the last item in the filenames. Because the loop will finish loop, the fname variable will keep being override util the loop is finish
		// By the time the go routine actually run, each of them will go through the same varibale "fname".
		// The behavior is fix in Go 1.22 where each time we go through the loop, a new variable is created to whole the value.
		// https://go.dev/doc/faq#closures_and_goroutines
		// https://go.dev/blog/loopvar-preview
		go func(fname string) {
			defer wg.Done()
			f, err := os.Open(fname)
			if err != nil {
				errCh <- fmt.Errorf("Cannot open file: %w", err)
				return
			}

			// Parse the CSV into a slice of float64 numbers
			data, err := csv2float(f, column)
			if err != nil {
				errCh <- err
			}

			if err := f.Close(); err != nil {
				errCh <- err
			}

			resCh <- data

		}(fname)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
}
