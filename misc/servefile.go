package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", "localhost:11111", "http server address")
	flag.Parse()

	if flag.NArg() == 0 {
		usage("no file is provided")
	}

	files := NewFiles(flag.Args())
	http.Handle("/", http.FileServer(files))
	http.ListenAndServe(*addr, nil)
}

func usage(format string, args ...interface{}) {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, format+"\n", args...)
	fmt.Fprintf(out, "Usage: %s <file> ...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

type Files struct {
	files map[string]struct{}
}

func NewFiles(list []string) Files {
	fs := Files{make(map[string]struct{})}
	for _, f := range list {
		f = filepath.Clean(f)
		if filepath.IsAbs(f) {
			usage("Absolute path:%q is not supported", f)
		}
		if _, err := os.Stat(f); err != nil {
			usage("Stat file:%q err:%v", f, err)
		}
		fs.files[f] = struct{}{}
	}
	return fs
}

func (fs Files) Open(name string) (http.File, error) {
	name = name[1:]
	if _, ok := fs.files[name]; !ok {
		return nil, os.ErrNotExist
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
