package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// TODO:
// impl double check(binary compare) for conflicted hashes but not same contents
// define type struct for JSON
// consider API: Load, Save

// base information
var (
	Name    = filepath.Base(os.Args[0])
	Version = "1.0.0dev"
)

var usageWriter io.Writer = os.Stderr

func makeUsage(w *io.Writer) func() {
	return func() {
		flag.CommandLine.SetOutput(*w)
		// two spaces
		fmt.Fprintf(*w, "Description:\n")
		fmt.Fprintf(*w, "  Find duplicate files\n\n")
		fmt.Fprintf(*w, "Usage:\n")
		fmt.Fprintf(*w, "  %s [Options] -- [PATH]...\n\n", Name)
		fmt.Fprintf(*w, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(*w, "\n")
		examples := `Examples:
  $ ` + Name + ` -help                         # Display help message
  $ ` + Name + ` /path/file1 /path/file2       # Compare two files
  $ ` + Name + ` /path/dir                     # Check recursive
  $ ` + Name + ` -verbose /path/dir /path/file # With verbose
`
		fmt.Fprintf(*w, "%s\n", examples)
	}
}

var (
	Log    *log.Logger = log.New(ioutil.Discard, Name+":", log.LstdFlags)
	ErrLog *log.Logger = log.New(os.Stderr, Name+":[err]:", log.LstdFlags)
)

var opt struct {
	version bool
	help    bool

	verbose bool
	hash    string

	paths string

	// TODO: impl double check, compare binaries
	//double bool
}

func init() {
	flag.Usage = makeUsage(&usageWriter)

	flag.BoolVar(&opt.version, "version", false, "Display version")
	flag.BoolVar(&opt.help, "help", false, "Display help message")

	flag.BoolVar(&opt.verbose, "verbose", false, "Verbose output")
	flag.StringVar(&opt.hash, "hash", "sha256", "Specify hash function")

	flag.StringVar(&opt.paths, "paths", "", "Specify search paths (separator is \""+string(filepath.ListSeparator)+"\")")

	flag.Parse()
}

// TODO:
// 1. impl struct and separate to methods
//type DuplicateFiles struct {
//	// hash function
//	Function string              `json:"function"`
//	HashMap  map[string][]string `json:"hash_map"`
//}
// 2. remove? stdout and stderr
// 3. consider to use goroutine for calculate of hash
func run(stdout io.Writer, stderr io.Writer, hashF string, paths []string) error {
	var hash hash.Hash
	switch hashF {
	case "sha256":
		hash = sha256.New()
	case "sha1":
		hash = sha1.New()
	case "md5":
		hash = md5.New()
	default:
		return fmt.Errorf("invalid hash function: %q", hashF)
	}

	type dup struct {
		isDupl bool
		sum    []byte
		path   []string
	}
	hashMap := make(map[string]*dup)
	avoidMap := make(map[string]bool)

	for _, root := range paths {
		abs, err := filepath.Abs(root)
		if err != nil {
			ErrLog.Println(err)
			continue
		}
		err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				switch {
				case os.IsPermission(err) || os.IsNotExist(err):
					// TODO: set exit code to not 0
					ErrLog.Println(err)
					return nil
				default:
					return err
				}
			}
			if !info.Mode().IsRegular() || avoidMap[path] {
				return nil
			}
			avoidMap[path] = true

			Log.Printf("checking: %q\n", path)

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			hash.Reset()
			if _, err := io.Copy(hash, f); err != nil {
				return err
			}
			key := fmt.Sprintf("%x", hash.Sum(nil))
			if d, ok := hashMap[key]; ok {
				d.isDupl = true
				d.path = append(d.path, path)
			} else {
				hashMap[key] = &dup{isDupl: false, sum: hash.Sum(nil), path: []string{path}}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(stdout, "Used hash function: %q\n", hashF)
	for _, d := range hashMap {
		if d.isDupl {
			fmt.Fprintf(stdout, "Conflicted hash [%x]\n", d.sum)
			for _, s := range d.path {
				fmt.Fprintf(stdout, "\t%s\n", s)
			}
		}
	}

	return nil
}

func main() {
	switch {
	case opt.version:
		fmt.Printf("%s %s\n", Name, Version)
		os.Exit(0)
	case opt.help:
		usageWriter = os.Stdout
		flag.Usage()
		os.Exit(0)
	}

	if opt.verbose {
		Log.SetOutput(os.Stdout)
	}

	var paths []string
	if opt.paths != "" {
		paths = append(paths, filepath.SplitList(opt.paths)...)
	}
	paths = append(paths, flag.Args()...)

	if err := run(os.Stdout, os.Stderr, opt.hash, paths); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
