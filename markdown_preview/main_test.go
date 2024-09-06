package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenfile = "./testdata/test1.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)

	expected, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Errorf("Result parseContent does not match golden file:\n%s", cmp.Diff(expected, result))
	}

}

func TestRun(t *testing.T) {
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Errorf("Result content does not match golden file:\n%s", cmp.Diff(expected, result))
	}

	os.Remove(resultFile)
}
