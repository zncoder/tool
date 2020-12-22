package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zncoder/assert"
)

func main() {
	varName := flag.String("v", "", "template var name, guess from template file name if not set")
	pkgName := flag.String("p", "", "package name, guess from current dir if not set")
	file := flag.String("f", "", "template file, name must be valid go identifier")
	flag.Parse()

	if *varName == "" {
		*varName = setVarName(*file)
	}
	if *pkgName == "" {
		*pkgName = guessPkgName()
	}

	out, err := os.Create(*file + ".go")
	assert.Nil(err)
	w := bufio.NewWriter(out)
	defer out.Close()
	defer w.Flush()

	fmt.Fprintf(w, `package %s

import "html/template"

var %s = template.Must(template.New("%s").Parse(%s))
`,
		*pkgName, *varName, *varName, formatFile(*file))
}

func formatFile(fn string) string {
	f, err := os.Open(fn)
	assert.Nil(err)
	defer f.Close()

	buf := &bytes.Buffer{}

	var inBq bool
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ln := sc.Text()
		if strconv.CanBackquote(ln) {
			if !inBq {
				buf.WriteByte('`')
				inBq = true
			}
			buf.WriteString(ln)
			buf.WriteByte('\n')
		} else {
			if inBq {
				buf.WriteString("` + \n")
				inBq = false
			}
			buf.WriteByte('\t')
			buf.WriteString(strconv.Quote(ln + "\n"))
			buf.WriteString(" + \n")
		}
	}
	if inBq {
		buf.WriteByte('`')
	} else {
		buf.WriteString(`""`)
	}
	return buf.String()
}

func setVarName(fn string) string {
	fn = filepath.Base(fn)
	i := strings.Index(fn, ".")
	if i < 0 {
		return fn
	}
	return fn[:i]
}

func guessPkgName() string {
	matches, err := filepath.Glob("*.go")
	assert.Nil(err)
	assert.OKf(len(matches) > 0, "no go file is found to guess package")

	f, err := os.Open(matches[0])
	assert.Nil(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ln := sc.Text()
		if !strings.HasPrefix(ln, "package ") {
			continue
		}
		pkg := strings.TrimSpace(strings.TrimPrefix(ln, "package "))
		assert.OKf(pkg != "", "empty package name")
		return pkg
	}
	assert.OKf(false, "no package name")
	return ""
}
