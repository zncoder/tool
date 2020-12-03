package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mdp/qrterminal"
	"github.com/zncoder/assert"
	"github.com/zncoder/easycmd"
)

func main() {
	easycmd.Handle("download", runDownload, "embed the download page in QR")
	easycmd.Handle("snippet", runSnippet, "embed the snippet in QR")
	easycmd.Main()
}

func runDownload() {
	port := flag.Int("p", 10030, "port")
	flag.Parse()
	assert.OK(flag.NArg() > 0)

	type FileInfo struct {
		Name string
		Path string
	}
	var files []FileInfo
	for _, f := range flag.Args() {
		files = append(files, FileInfo{Name: filepath.Base(f), Path: f})
	}

	status := make(chan string, 1)

	var wg sync.WaitGroup
	wg.Add(len(files))

	tmpl := template.Must(template.New("index").Parse(indexTmpl))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := tmpl.Execute(w, files)
		assert.Nil(err)
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		log.Println("req", r)
		filename := r.URL.Query().Get("f")
		f, err := os.Open(filename)
		if err != nil {
			http.Error(w, "open file", http.StatusNotFound)
			return
		}
		defer f.Close()

		st, err := f.Stat()
		assert.Nil(err)
		log.Printf("Downloading file:%q size:%s", filename, printSize(st.Size()))

		w.Header().Add("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, f)
		if err != nil {
			log.Printf("download file:%q err:%v", filename, err)
			return
		}
		status <- filename
		wg.Done()
		log.Println("finish", filename)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		x := <-status
		log.Println("status:", x)
		io.WriteString(w, x)
	})

	ip := HostIP()
	lr, err := net.Listen("tcp", fmt.Sprintf("%v:%d", ip, *port))
	assert.Nil(err)
	go http.Serve(lr, nil)

	addr := fmt.Sprintf("http://%v/", lr.Addr())
	log.Printf("Listening on %s", addr)
	//qrterminal.GenerateHalfBlock(addr, qrterminal.H, os.Stdout)
	qrterminal.Generate(addr, qrterminal.M, os.Stdout)
	wg.Wait()
	lr.Close()
	time.Sleep(10 * time.Second)
}

func printSize(n int64) string {
	switch {
	case n < 1e6:
		return strconv.FormatInt(n, 10)
	case n < 1e9:
		return fmt.Sprintf("%.2fMB", float64(n)/1e6)
	default:
		return fmt.Sprintf("%.2fGB", float64(n)/1e9)
	}
}

func HostIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	assert.Nil(err)
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

func runSnippet() {
	flag.Parse()
	assert.OK(flag.NArg() > 0)

	s := strings.Join(flag.Args(), " ")
	qrterminal.GenerateHalfBlock(s, qrterminal.M, os.Stdout)
}
