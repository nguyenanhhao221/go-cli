package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1, word2, word3, word4\n")

	exp := 4
	opts := wcOpts{lines: false}
	res := count(b, &opts)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1, word2, word3, word4\nline2\nline3 word1")

	exp := 3
	opts := wcOpts{lines: true}
	res := count(b, &opts)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("gopher")

	exp := 6
	opts := wcOpts{countBytes: true}
	res := count(b, &opts)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}
