package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// gc_Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func gc_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gc_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gc_endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(gc_suiteRE).FindStringSubmatch

	var suites = make([]*Suite, 0)
	var suiteName string
	var curSuite *Suite

	var testName string
	var message []string

	scanner := bufio.NewScanner(rd)
	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		tokens := find_start(line)
		if len(tokens) > 0 {
			if testName != "" {
				return nil, fmt.Errorf("%d: start in middle\n", lnum)
			}
			suiteName = tokens[1]
			testName = tokens[2]
			message = []string{}
			continue
		}

		tokens = find_end(line)
		if len(tokens) > 0 {
			if testName == "" {
				return nil, fmt.Errorf("%d: orphan end", lnum)
			}
			if (tokens[2] != suiteName) || (tokens[3] != testName) {
				return nil, fmt.Errorf("%d: suite/name mismatch", lnum)
			}
			curTest := &Test{Name: testName}
			curTest.Message = strings.Join(message, "\n")
			curTest.Time = tokens[4]
			curTest.Failed = (tokens[1] == "FAIL") || (tokens[1] == "PANIC")
			curTest.Passed = (tokens[1] == "PASS")
			curTest.Skipped = (tokens[1] == "SKIP")

			if curSuite == nil || curSuite.Name != suiteName {
				curSuite = &Suite{Name: suiteName}
				suites = append(suites, curSuite)
			}
			curSuite.Tests = append(curSuite.Tests, curTest)

			testName = ""
			suiteName = ""
			message = []string{}

			continue
		}

		// last "suite" is test summary
		tokens = find_suite(line)
		if tokens != nil {
			if curSuite == nil {
				curSuite = &Suite{Name: tokens[2], Status: tokens[1], Time: tokens[3]}
				suites = append(suites, curSuite)
			} else {
				curSuite.Status = tokens[1]
				curSuite.Time = tokens[3]
			}

			testName = ""
			suiteName = ""
			message = []string{}

			continue
		}

		if testName != "" {
			message = append(message, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}
