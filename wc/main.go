package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	fmt.Println(count(os.Stdin))
}

func count(r io.Reader) int {
	// Create a new scanner to prepare to read from io.Reader
	scanner := bufio.NewScanner(r)

	// Define the scanner split type to words (default is split by lines)
	scanner.Split(bufio.ScanWords)

	// Defining a counter
	var wc int = 0

	for scanner.Scan() {
		wc++
	}

	return wc

}
