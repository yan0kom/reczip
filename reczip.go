package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var rootPath string
var extension string
var deleteSource bool

var filterByExtension = false
var filenames = make(chan string, 100)

func init() {
	flag.StringVar(&rootPath, "path", "", "path to files")
	flag.StringVar(&extension, "ext", "", "extension of files to archive")
	flag.BoolVar(&deleteSource, "del", false, "delete source file")
}

func main() {
	flag.Parse()
	if rootPath == "" || extension == "" || flag.NFlag() != 3 {
		flag.Usage()
		os.Exit(1)
	}


	filterByExtension = len(extension) > 0

	fmt.Printf("Input path: %s\n", rootPath)
	if filterByExtension {
		fmt.Printf("Filter by extension: %s\n", extension)
	}

	go enumerateFilenames()

	for fn := range filenames {
		fmt.Printf("%s\n", fn)
		if err := zipFile(fn); err != nil {
			log.Fatal(err)
		}
		if deleteSource {
			if err := os.Remove(fn); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func enumerateFilenames() {
	defer close(filenames)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if !filterByExtension || filepath.Ext(path) == extension {
				filenames <- path
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func zipFile(fn string) error {
	zfn := fn + ".zip"
	zf, err := os.Create(zfn)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	fr, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fr.Close()

	zfw, err := zw.Create(filepath.Base(fn))
	if err != nil {
		return err
	}
	_, err = io.Copy(zfw, fr)
	if err != nil {
		return err
	}

	return nil
}
