package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	/// init
	// test directory
	testRoot, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	testRoot = filepath.Join(testRoot, "main")
	if err := os.MkdirAll(testRoot, 0777); err != nil {
		t.Fatal(err)
	}
	// test files
	var (
		sameFilesContents = "hello world"
		sameFiles         = []string{
			filepath.Join(testRoot, "same_one.txt"),
			filepath.Join(testRoot, "same_two.txt"),
		}

		uniqueFilesContents = "unique"
		uniqueFile          = filepath.Join(testRoot, "unique.txt")
	)

	// write
	write := func(path string, b []byte) {
		if err := ioutil.WriteFile(path, b, 0666); err != nil {
			t.Fatal(err)
		}
	}
	write(sameFiles[0], []byte(sameFilesContents))
	write(sameFiles[1], []byte(sameFilesContents))
	write(uniqueFile, []byte(uniqueFilesContents))

	// output buffer
	var (
		buf    = bytes.NewBuffer([]byte{})
		errbuf = bytes.NewBuffer([]byte{})
	)

	// TODO: append fatal case
	if exit := Run(buf, errbuf, "sha512_256", true, []string{testRoot}); exit != 0 {
		t.Fatal(errbuf)
	}

	// TODO: consider
	if buf.Len() != 0 {
		t.Log("buf:", buf.String())
	}
	if errbuf.Len() != 0 {
		t.Log("errbuf:", errbuf.String())
	}
}
