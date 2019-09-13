package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type Response struct {
	*http.Response
	err error
}

const (
	// Prod
	//  sctevsRaw = "https://sctevs.linear.theplatform.com/scte224/raw"
	// Stage
	sctevsRaw = "https://sctevs.linear.preview.theplatform.com/scte224/raw"
)

var (
	concurrentMax int
	Account       string
	Token         string
)

func init() {
	flag.IntVar(&concurrentMax, "concurrent", 200, "Maximum concurrent requests")
	flag.StringVar(&Account, "account", "", "Account number")
	flag.StringVar(&Token, "token", "", "MPX Token")
}

func main() {

	flag.Parse()

	if Token == "" || Account == "" {
		fmt.Println("Token and account is required.")
		return
	}

	urls := []string{}

	// Open the file.
	f, _ := os.Open("ids.txt")

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		url := sctevsRaw + "/" + Account + "/" + strings.TrimSpace(line) + "?token=" + Token
		urls = append(urls, url)
	}

	if len(urls) == 0 {
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	reqChan := make(chan *http.Request)
	respChan := make(chan Response)

	start := time.Now()
	go dispatcher(urls, reqChan)
	go workerPool(reqChan, respChan)
	conns, size := consumer(len(urls), respChan)
	took := time.Since(start)
	ns := took.Nanoseconds()
	av := ns / conns
	average, err := time.ParseDuration(fmt.Sprintf("%d", av) + "ns")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Connections:\t%d\nConcurrent:\t%d\nTotal size:\t%d bytes\nTotal time:\t%s\nAverage time:\t%s\n", conns, concurrentMax, size, took, average)

}

// Dispatcher
func dispatcher(urls []string, reqChan chan *http.Request) {
	defer close(reqChan)
	for _, url := range urls {
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			log.Println(err)
		}
		reqChan <- req
	}
}

// Worker Pool
func workerPool(reqChan chan *http.Request, respChan chan Response) {
	t := &http.Transport{}
	for i := 0; i < concurrentMax; i++ {
		go worker(t, reqChan, respChan)
	}
}

// Worker
func worker(t *http.Transport, reqChan chan *http.Request, respChan chan Response) {
	for req := range reqChan {
		resp, err := t.RoundTrip(req)
		r := Response{resp, err}
		respChan <- r
	}
}

// Consumer
func consumer(urlCount int, respChan chan Response) (int64, int64) {
	var (
		conns int64
		size  int64
	)

	for conns < int64(urlCount) {
		select {
		case r, ok := <-respChan:
			if ok {
				if r.err != nil {
					log.Println(r.err)
				} else {
					size += r.ContentLength
					if err := r.Body.Close(); err != nil {
						log.Println(r.err)
					}
				}
				conns++
			}
		}
	}
	return conns, size
}
