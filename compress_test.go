package main

import (
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
