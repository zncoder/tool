package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	parseLayout  = "20060102"
	outputLayout = "2006-01-02T15:04:05.000MST"
)

func main() {
	showMillis := flag.Bool("m", false, "show in millisecond")
	flag.Parse()

	arg := flag.Arg(0)
	if len(arg) == len(parseLayout) {
		toSecond(arg, *showMillis)
	} else {
		toDate(arg)
	}
}

func toSecond(arg string, millis bool) {
	t, err := time.Parse(parseLayout, arg)
	if err != nil {
		log.Fatalf("parse %q", arg)
	}
	if millis {
		fmt.Println(t.UnixNano() / 1e6)
	} else {
		fmt.Println(t.Unix())
	}
}

func toDate(arg string) {
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		log.Fatalf("parse second %q", arg)
	}
	now := time.Now()
	var t time.Time
	if n > 100*now.Unix() {
		// millis
		t = time.Unix(n/1000, n%1000)
	} else {
		t = time.Unix(n, 0)
	}
	fmt.Printf("%s    %s\n", t.In(time.UTC).Format(outputLayout), t.In(time.Local).Format(outputLayout))
}
