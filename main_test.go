package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenLogFile_CreatesParentDir(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "logs", "site.log")

	f, err := openLogFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
	if _, err := f.WriteString("ok\n"); err != nil {
		t.Fatal(err)
	}
}
