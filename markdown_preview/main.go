package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	if err := run(*filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run Read all data from input file, call parseContent and saveHtml
func run(filename string, out io.Writer) error {
	// Read all data from input file and check for error
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)
	// Create a temp file to save the html later
	tmp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	outName := tmp.Name()
	//NOTE: Here we write to an io.Writer, then later in the test we can read from this.
	// This approach is good, since we don't need to directly return the tmp file name in the actual function. Only we need to know the tmp file when we test to compare the result.
	fmt.Fprintln(out, outName)
	if err := saveHtml(outName, htmlData); err != nil {
		return err
	}

	return preview(outName)

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

func preview(filename string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "darwin":
		cName = "open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	default:
		return fmt.Errorf("OS is not supported")
	}

	cParams = append(cParams, filename)

	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	return exec.Command(cPath, cParams...).Run()
}
