package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	// extension to filter out
	ext string
	// Min file size
	size int64
	// list files
	list bool
	// delete files
	del bool
	// log destination writer
	wLog    io.Writer
	archive string
}

func main() {
	root := flag.String("root", "", "Root directory to start")
	list := flag.Bool("list", false, "List file only")
	//Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	del := flag.Bool("del", false, "Delete files")
	logFile := flag.String("log", "", "Log deletes to this file")
	archive := flag.String("archive", "", "Archive directory")
	flag.Parse()

	var f = os.Stdout

	// If user specify logFile we set f which represent a io.Writer interface into this file to be prepare to write into it
	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			f = file
		}
		defer f.Close()

	}

	if *archive != "" {
		f, err := os.Stat(*archive)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !f.IsDir() {
			fmt.Fprintf(os.Stderr, "Error:%s is not a directory\n", *archive)
			os.Exit(1)
		}
	}

	c := config{
		ext:     *ext,
		list:    *list,
		size:    *size,
		del:     *del,
		archive: *archive,
		wLog:    f,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}
		if cfg.list {
			return listFile(path, out)
		}

		// Archive file
		if cfg.archive != "" {
			if err := archiveFile(path, root, cfg.archive); err != nil {
				return err
			}
		}

		// Delete file
		if cfg.del {
			return delFile(path, delLogger)
		}

		return listFile(path, out)
	})

}
