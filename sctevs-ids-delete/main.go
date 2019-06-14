package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	// Prod
	//  sctevsRaw = "https://sctevs.linear.theplatform.com/scte224/raw"
	// Stage
	sctevsRaw = "https://sctevs.linear.preview.theplatform.com/scte224/raw"
)

var (
	Account string
	Token   string
)

func main() {
	flag.StringVar(&Account, "account", "", "Account number")
	flag.StringVar(&Token, "token", "", "MPX Token")
	flag.Parse()

	if Token == "" || Account == "" {
		fmt.Println("Token and account is required.")
		return
	}

	list := []string{}
	// Open the file.
	f, _ := os.Open("ids.txt")
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		list = append(list, strings.TrimSpace(line))
	}

	for _, id := range list {
		url := sctevsRaw + "/" + Account + "/" + id + "?token=" + Token
		fmt.Println(url)
		request(url)
	}
}

func request(url string) error {
	var httpClient http.Client
	var request *http.Request
	var contentBuffer io.Reader
	var err error

	request, err = http.NewRequest(http.MethodDelete, url, contentBuffer)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to create request: %v", err)
		fmt.Println(errorMsg)
	}

	_, err = httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
