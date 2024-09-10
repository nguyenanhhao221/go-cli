package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		expected string
		cfg      config
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 20, list: true}, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead \n", tc.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		extNoDelete string
		expected    string
		cfg         config
		nDelete     int
		nNodelete   int
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNodelete: 10, expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNodelete: 0, expected: ""},
		{name: "DeleteExtensionMix", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNodelete: 5, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer, logBuffer bytes.Buffer

			files := map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNodelete,
			}

			// Remember because this test is mainly test the run function, we cannot pass the -del log file name like in the cli
			// There for we use this logBuffer to let our test run write to this logBuffer, then we will validate the output against it
			tc.cfg.wLog = &logBuffer

			rootTmpDir := createTempDirWithMockFiles(t, files)

			err := run(rootTmpDir, &buffer, tc.cfg)
			if err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead \n", tc.expected, res)
			}

			fileLeft, err := os.ReadDir(rootTmpDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(fileLeft) != tc.nNodelete {
				t.Errorf("Expected %d, got %d instead \n", tc.nNodelete, len(fileLeft))
			}

			expLogLines := tc.nDelete + 1
			lines := len(bytes.Split(logBuffer.Bytes(), []byte("\n")))
			if lines != expLogLines {
				t.Errorf("Expected %d log lines, got %d", expLogLines, lines)
			}
		})

	}

}

func createTempDirWithMockFiles(t *testing.T, files map[string]int) string {
	t.Helper()
	tempDir := t.TempDir()
	for k, v := range files {
		for i := 0; i < v; i++ {
			fname := fmt.Sprintf("file*%s", k)
			f, err := os.CreateTemp(tempDir, fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if _, err := f.Write([]byte("dummy")); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir
}
