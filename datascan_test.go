package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

var reg = regexp.MustCompile(expression)

func TestRegExp(t *testing.T) {

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
		capt := reg.FindStringSubmatch(fds.Text())
		if len(capt) != 0 {
			fmt.Println(capt[2])
			fmt.Println(capt[1])
			ts, err := time.Parse(layout, capt[1])
			if err != nil {
				fmt.Printf("Could not parse timestamp %s, %s\n", capt[1], err)
			}
			fmt.Println(ts.String())
		}
	}
}

func TestNewRecordScanner(t *testing.T) {
	rs, err := newRecordScanner(path + "/testdata/filtertest")
	if err != nil {
		fmt.Println(err)
	}

	for rs.ScanDedup() {
		fmt.Printf("%s, %s, %s\n", rs.Data().Time, rs.Data().Name, rs.Error())
	}
}
