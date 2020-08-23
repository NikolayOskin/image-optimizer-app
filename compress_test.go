package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestByteCountSI(t *testing.T) {
	cases := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1000, "1.0 kB"},
		{1500, "1.5 kB"},
		{1000000, "1.0 MB"},
		{10000000, "10.0 MB"},
	}

	for _, testcase := range cases {
		result := ByteCountSI(testcase.bytes)

		if testcase.expected != result {
			t.Errorf("returned wrong string: got %s want %s", result, testcase.expected)
		}
	}
}

func TestCountCompressedSizes(t *testing.T) {
	cases := []struct {
		inputFileSize           int64
		outFileSize             int64
		expectSavedAmountIsZero bool
		expectError             bool
	}{
		{0, 0, true, true},
		{1000, 1500, true, false},
		{1000, 500, false, false},
		{1000, 1000, true, false},
	}

	for _, testcase := range cases {
		result := compressResult{}

		tmpfile := filepath.Join(imagesPath, "tmpfile")
		if err := ioutil.WriteFile(
			tmpfile,
			make([]byte, testcase.outFileSize),
			0666,
		); err != nil {
			t.Errorf("can't create temp file: %v", err)
		}
		file, err := os.Open(tmpfile)
		if err != nil {
			t.Errorf("can't open tmpfile: %v", err)
		}

		err = result.countCompressedSizes(testcase.inputFileSize, file)
		if err == nil && testcase.expectError {
			t.Error("expected error, getting nil")
		}
		if err != nil && testcase.expectError == false {
			t.Error("error is returned, but not expected")
		}
		if testcase.expectSavedAmountIsZero && result.saved != 0 {
			t.Errorf("expected saved is 0, %v returned", result.saved)
		}
		if result.saved == 0 && !testcase.expectSavedAmountIsZero {
			t.Error("saved is 0, expected is more than 0")
		}
	}
}
