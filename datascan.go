package main

import (
	"regexp"
	"strings"
	"time"
)

const (
	// regular expression to capture timestamp and membername from loglines
	// containing 'welcome' timestamp is submatch imdex 1 and membername
	// submatch index 2
	expression = `\[([[:graph:]]*).*\](?:.*welcome\/)([[:alnum:]]*)`
	// layout tp parse timestamp
	layout = "15/Jan/2006:15:04:05"
)

type record struct {
	time.Time
	Name string
}

type recordScanner struct {
	*logDirFilter
	*regexp.Regexp
	Rec record
	err error
}

func (rs *recordScanner) Scan() bool {
	// scan progresses to next filtered line in log dir
	if rs.logDirFilter.Scan() {

		// capture timestamp and membername from line
		capture := rs.Regexp.FindStringSubmatch(rs.logDirFilter.Text())
		// parse timestamp
		stamp, err := time.Parse(layout, capture[1])
		// parse membername
		name := capture[2]

		// create and assign new record
		(*rs).Rec = record{stamp, name}
		// assign error value
		(*rs).err = err
		return true
	}
	return false
}
func (rs *recordScanner) ScanDedup() bool {
	last := rs.Rec.Name
	for (*rs).Scan() {
		if (*rs).Rec.Name != last {
			return true
		}
	}
	return false
}

func (rs *recordScanner) Data() record {
	return rs.Rec
}

func (rs *recordScanner) Error() error {
	return rs.err
}
func newRecordScanner(path string) (rs *recordScanner, err error) {
	// create log dir line scanner
	lds, err := openLogDir(path)
	if err != nil {
		return rs, err
	}
	// create log dir line filter for all lines containing the string 'welcome'
	ldf := newLogDirFilter(
		lds,
		func(line string) bool {
			return strings.Contains(line, "welcome")
		},
	)

	reg := regexp.MustCompile(expression)

	rs = &recordScanner{
		ldf,
		reg,
		record{
			time.Time{},
			"",
		},
		nil,
	}

	return rs, err
}
