package cleaner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDeleteOldFiles(t *testing.T) {
	// creating temp dir and file inside it
	content := []byte("some bytes")
	dir, _ := ioutil.TempDir("./", "example")
	defer os.RemoveAll(dir)

	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Error("can't write file to temp dir")
	}

	time.Sleep(15 * time.Nanosecond)

	cases := []struct {
		olderThan time.Duration
		expected  bool
	}{
		{1 * time.Hour, false},
		{1 * time.Minute, false},
		{10 * time.Second, false},
		{1 * time.Nanosecond, true},
	}

	for _, testcase := range cases {
		err := deleteOldFiles(dir, testcase.olderThan)
		if err != nil {
			t.Error("error while deleting file from dir")
		}

		_, err = os.Stat(tmpfn)
		if os.IsNotExist(err) != testcase.expected {
			t.Error("file must be deleted but it exists")
		}
	}
}
