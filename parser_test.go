package main

import (
	"reflect"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func Test_Parser(t *testing.T) {

	var tests = []struct {
		file string
		test *Test
	}{
		// Single field statement
		{
			file: `gotest-pass.out`,
			test: &Test{
				Name:   "TestAdd",
				Passed: true,
				Time:   "0.01",
			},
		},
	}

	for i, tt := range tests {
		rd := getInputData(tt.file)

		test, err := NewParser(rd).Parse()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(tt.test, test) {
			t.Errorf("%d. %q\nexp=%#v\ngot=%#v\n\n", i, tt.file, tt.test, test)
		}
	}
}
