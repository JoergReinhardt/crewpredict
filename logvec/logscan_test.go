package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var path = os.Getenv("GOPATH") + "/src/github.com/JoergReinhardt/crewpredict"
var lds *logDirScanner
var err error

func TestScannerFromFile(t *testing.T) {

	f, err := os.Open(path + "/testdata/dirtest/1.log")
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
	}
	fmt.Printf("file: %s\n", f)

	scan, err := scannerFromFile(f)
	if err != nil {
		fmt.Printf("error creating scanner: %s\n", err)
	}
	for scan.Scan() {
		fmt.Printf("scanner: %s\n", scan.Text())
	}

	f, err = os.Open(path + "/testdata/dirtest/2.log.gz")
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
	}
	fmt.Printf("file: %s\n", f)

	scan, err = scannerFromFile(f)
	if err != nil {
		fmt.Printf("error creating scanner: %s\n", err)
	}
	for scan.Scan() {
		fmt.Printf("scanner: %s\n", scan.Text())
	}
}

func TestOpenLogDirScanner(t *testing.T) {

	lds, err = openLogDir(path + "/testdata/dirtest")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("list: %s\n", lds.list)
	fmt.Printf("idx: %d\n", lds.idx)
	fmt.Printf("error: %s\n", lds.Err())

	for lds.Scan() {
		fmt.Printf("lds len, idx, text: %d, %d, %s\n", len(lds.list), lds.idx, lds.Text())
	}
}

func TestNewLogDirFilter(t *testing.T) {
	lds, err = openLogDir(path + "/testdata/filtertest")
	if err != nil {
		fmt.Println(err)
		return
	}
	fds := newLogDirFilter(
		lds,
		func(line string) bool {
			return strings.Contains(line, "welcome")
		},
	)

	for fds.Scan() {
		fmt.Println(fds.Text())
	}
}
