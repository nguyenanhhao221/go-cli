package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>Markdown Preview Tool</title>
</head>
<body>
	`
	footer = `
</body>
</html>
`
)

func main() {
	// Parse flag
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	//
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run Read all data from input file, call parseContent and saveHtml
func run(filename string) error {
	// Read all data from input file and check for error
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	return saveHtml(outName, htmlData)

}

func parseContent(input []byte) []byte {
	// Parse the input markdown file through blackfriday and bluemonday to get sanitize html
	out := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(out)

	var buffer bytes.Buffer

	// Write html to buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()

}

func saveHtml(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}
