package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {

	// Defining  a boolean flag -l count lints instead of words
	var lines *bool = flag.Bool("l", false, "Count lines")
	// Parsing the flag provided by user
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines))
}

func count(r io.Reader, isCountLines bool) int {
	// Create a new scanner to prepare to read from io.Reader
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, we want to count words
	if !isCountLines {
		// Define the scanner split type to words (default is split by lines)
		scanner.Split(bufio.ScanWords)
	}

	// Defining a counter
	var wc int = 0

	for scanner.Scan() {
		wc++
	}

	return wc

}
