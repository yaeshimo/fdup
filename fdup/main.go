package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// base information
const (
	Name                 = "fdup"
	Version              = "1.0.0dev"
	DefaultHashAlgorithm = "sha512_256"
)

type option struct {
	version bool
	verbose bool
	hash    string
}

var opt = &option{}

func init() {
	flag.BoolVar(&opt.version, "version", false, "show version")
	flag.BoolVar(&opt.verbose, "verbose", false, "verbose")
	flag.StringVar(&opt.hash, "hash", DefaultHashAlgorithm, "specify use hash algorithmm")
}

// Run main
func Run(stdout, stderr io.Writer, usehash string, verbose bool, targets []string) int {
	log.SetPrefix("[" + Name + "]:")
	log.SetOutput(stdout)
	errlogger := log.New(stderr, "["+Name+"]:", log.Lshortfile)

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	var hash hash.Hash
	switch usehash {
	//case "md5":
	//	hash = md5.New()
	case "sha512_256":
		hash = sha512.New512_256()
	default:
		fmt.Fprintln(stderr, "invalid hash algorithm:", usehash)
		return 1
	}

	type dup struct {
		isDupl bool
		sum    []byte
		path   []string
	}
	// key=string(hash.Sum(nil))
	hashMap := make(map[string]*dup)
	// key=fullFilePath for avoid duplicate check
	avoidMap := make(map[string]bool)
	for _, root := range targets {
		abs, err := filepath.Abs(root)
		if err != nil {
			errlogger.Println(err)
			continue
		}
		err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				if avoidMap[path] {
					return nil
				}
				avoidMap[path] = true

				log.Println("check:", path)

				f, err := os.Open(path)
				if err != nil {
					errlogger.Println(err)
					return nil
				}
				defer f.Close()

				hash.Reset()
				if _, err := io.Copy(hash, f); err != nil {
					errlogger.Println(err)
					return nil
				}
				key := string(hash.Sum(nil))
				if d, ok := hashMap[key]; ok {
					d.isDupl = true
					d.path = append(d.path, path)
				} else {
					hashMap[key] = &dup{isDupl: false, sum: hash.Sum(nil), path: []string{path}}
				}
			}
			return nil
		})
		if err != nil {
			errlogger.Println(err)
		}
	}

	fmt.Fprintf(stdout, "Used hash algorithm: %q\n", usehash)
	for _, d := range hashMap {
		if d.isDupl {
			fmt.Fprintf(stdout, "Conflicted hash %v\n", d.sum)
			for _, s := range d.path {
				fmt.Fprintf(stdout, "\t%q\n", s)
			}
		}
	}
	return 0
}

func main() {
	flag.Parse()
	if opt.version {
		fmt.Fprintf(os.Stdout, "%s version %s\n", Name, Version)
		os.Exit(0)
	}
	os.Exit(Run(os.Stdout, os.Stderr, opt.hash, opt.verbose, flag.Args()))
}
