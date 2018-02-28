package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// TODO: impl
func TestRun(t *testing.T) {
	const helloWorld = "hello world"
	testDir, err := filepath.Abs("t")
	if err != nil {
		t.Fatal(err)
	}
	write := func(base string, b []byte) {
		if err := ioutil.WriteFile(filepath.Join(testDir, base), b, 0666); err != nil {
			t.Fatal(err)
		}
	}
	write("same1.txt", []byte(helloWorld))
	write("same2.txt", []byte(helloWorld))
	write("uniq.txt", []byte("uniq"))

	var (
		buf    = bytes.NewBuffer([]byte{})
		errbuf = bytes.NewBuffer([]byte{})
	)
	if exit := run(buf, errbuf, "sha512_256", true, []string{testDir}); exit != 0 {
		t.Fatal(errbuf)
	}
	if !strings.Contains(buf.String(), "hash=[") {
		t.Fatalf("can not found same content:\n[buf]=%v\n[errbuf]=%v\n", buf, errbuf)
	} else {
		t.Log(buf)
	}
}
