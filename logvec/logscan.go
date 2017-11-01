package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type logDirScanner struct {
	scan *bufio.Scanner
	path string
	list []string
	idx  int
	err  error
}

func (lds *logDirScanner) Scan() bool {
	// return immediately if current scanner has more lines
	if lds.scan.Scan() {
		return true
	}
	// at EOF of current scanner, if there are more files to scan, create scanner from  next file in list
	if lds.idx < len(lds.list)-1 {
		// update index
		(*lds).idx = lds.idx + 1

		// open next file in list
		f, err := os.Open(lds.path + "/" + lds.list[lds.idx])
		if err != nil {
			(*lds).err = err
			return false
		}

		// create new scanner from file and swap old
		(*lds).scan, err = scannerFromFile(f)
		if err != nil {
			(*lds).err = err
			return false
		}

		// perform scan on new scanner and return true, when successful
		if lds.scan.Scan() {
			return true
		}
	}
	// return false, when all scanners are depleted
	return false
}
func (lds *logDirScanner) Err() error {
	return lds.err
}
func (lds *logDirScanner) Text() string {
	return lds.scan.Text()
}
func (lds *logDirScanner) Bytes() []byte {
	return lds.scan.Bytes()
}

func openLogDir(path string) (lds *logDirScanner, err error) {

	dir, err := os.Open(path)
	if err != nil {
		return lds, err
	}
	defer dir.Close()

	finfo, err := dir.Stat()
	if err != nil {
		return lds, err
	}
	if finfo.IsDir() != true {
		return lds, fmt.Errorf("Not a directory")
	}

	names, err := dir.Readdirnames(0)
	if err != nil {
		return lds, fmt.Errorf("Could not read dirnames: %s", err)
	}
	if len(names) == 0 {
		return lds, fmt.Errorf("Directory is emty")
	}

	filename := path + "/" + names[0]
	f, err := os.Open(filename)
	if err != nil {
		return lds, fmt.Errorf("Could not open File: %s; %s", filename, err)
	}

	scan, err := scannerFromFile(f)
	if err != nil {
		return lds, fmt.Errorf("Could not create Scanner from File: %s; %s", f, err)
	}

	lds = &logDirScanner{
		scan,
		path,
		names,
		0,
		nil,
	}

	return lds, err
}

func scannerFromFile(reader io.Reader) (*bufio.Scanner, error) {

	var scanner *bufio.Scanner
	//create a bufio.Reader so we can 'peek' at the first few bytes
	bReader := bufio.NewReader(reader)

	testBytes, err := bReader.Peek(16) //read a few bytes without consuming
	if err != nil {
		return nil, err
	}
	//Detect if the content is gzipped
	contentType := http.DetectContentType(testBytes)

	//If we detect gzip, then make a gzip reader, then wrap it in a scanner
	if strings.Contains(contentType, "x-gzip") {
		gzipReader, err := gzip.NewReader(bReader)
		if err != nil {
			return nil, err
		}

		scanner = bufio.NewScanner(gzipReader)

	} else {
		//Not gzipped, just make a scanner based on the reader
		scanner = bufio.NewScanner(bReader)
	}

	return scanner, nil
}

type logDirFilter struct {
	scan   *logDirScanner
	filter func(string) bool
	line   string
}

func (ldf *logDirFilter) Scan() bool {
	for ldf.scan.Scan() {
		if ldf.filter(ldf.scan.Text()) {
			(*ldf).line = ldf.scan.Text()
			return true
		}
	}
	return false
}

func (ldf *logDirFilter) Text() string {
	return ldf.line
}
func (ldf *logDirFilter) Byte() []byte {
	return []byte(ldf.line)
}

func (ldf *logDirFilter) Err() error {
	return ldf.Err()
}

func newLogDirFilter(scan *logDirScanner, filter func(string) bool) *logDirFilter {
	return &logDirFilter{
		scan,
		filter,
		"",
	}
}
