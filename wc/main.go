package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type wcOpts struct {
	lines      bool
	countBytes bool
}

func main() {

	// Defining  a boolean flag -l count lints instead of words

	opts := wcOpts{}
	flag.BoolVar(&opts.lines, "l", false, "Count lines")
	flag.BoolVar(&opts.countBytes, "b", false, "Count Bytes")

	// Parsing the flag provided by user
	flag.Parse()
	files := flag.Args()

	if opts.countBytes && opts.lines {
		log.Fatal("-l and -b cannot be use together")
	}

	var totalCount int = 0

	if len(files) > 0 {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err:%v\n", err)
			}
			defer f.Close()
			totalCount += count(f, &opts)
		}
		fmt.Println(totalCount)
	} else {
		fmt.Println(count(os.Stdin, &opts))
	}
}

func count(r io.Reader, opts *wcOpts) int {
	// Create a new scanner to prepare to read from io.Reader
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, we want to count words
	if !opts.lines {
		// Define the scanner split type to words (default is split by lines)
		scanner.Split(bufio.ScanWords)
	}

	if opts.countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	// Defining a counter
	var wc int = 0

	for scanner.Scan() {
		wc++
	}

	return wc

}
