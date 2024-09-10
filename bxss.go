package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	BannerColor  = "\033[1;34m%s\033[0m\033[1;36m%s\033[0m"
	TextColor    = "\033[1;0m%s\033[1;32m%s\n\033[0m"
	InfoColor    = "\033[1;0m%s\033[1;35m%s\033[0m"
	NoticeColor  = "\033[1;0m%s\033[1;34m%s\n\033[0m"
	WarningColor = "\033[1;33m%s%s\033[0m"
	ErrorColor   = "\033[1;31m%s%s\033[0m"
	DebugColor   = "\033[0;36m%s%s\033[0m"
)

func main() {
	// Flag variables
	var c int
	var p string
	var h string
	var hf string
	var a bool
	var t bool

	// The flag / arguments
	flag.IntVar(&c, "concurrency", 30, "Set the concurrency")
	flag.StringVar(&h, "header", "", "Set a single custom header")
	flag.StringVar(&hf, "headerFile", "", "Path to file containing headers to test")
	flag.StringVar(&p, "payload", "", "the blind XSS payload")
	flag.BoolVar(&a, "appendMode", false, "Append the payload to the parameter")
	flag.BoolVar(&t, "parameters", false, "Test the parameters for blind xss")

	// Parse the arguments
	flag.Parse()

	// The banner
	fmt.Printf(BannerColor, `
	  ____               
	 |  _ \              
 	 | |_) |_  _____ ___ 
	 |  * <\ \/ / *_/ __|
	 | |_) |>  <\__ \__ \
	 |____//_/\_\___/___/
	                     
                    
	`, "-- Coded by @z0idsec -- \n")

	// Check if the bxss payload is set and either header or headerFile is provided
	if p == "" || (h == "" && hf == "") {
		flag.PrintDefaults()
		return
	}

	var headers []string
	if hf != "" {
		var err error
		headers, err = readHeadersFromFile(hf)
		if err != nil {
			fmt.Printf(ErrorColor, "Error reading header file: ", err.Error())
			return
		}
	} else {
		headers = []string{h}
	}

	fmt.Printf(NoticeColor, "\n[-] Please Be Patient for bxss\n ", "")

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processHeaders(headers, p, a, t)
		}()
	}
	wg.Wait()
}

func readHeadersFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var headers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		headers = append(headers, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}

func processHeaders(headers []string, payload string, appendMode bool, isParameters bool) {
	scanner := bufio.NewScanner(os.Stdin)
	client := &http.Client{Timeout: 3 * time.Second}

	for scanner.Scan() {
		link := scanner.Text()
		for _, header := range headers {
			testbxss(client, payload, link, header, appendMode, isParameters)
		}
	}
}

func testbxss(client *http.Client, payload string, link string, header string, appendMode bool, isParameters bool) {
	time.Sleep(500 * time.Microsecond)
	fmt.Println("")
	fmt.Printf(NoticeColor, "[+] \tHeader:  ", header)
	fmt.Printf(TextColor, "[+] \tPayload: ", payload)
	fmt.Println("")

	// Make GET Request
	makeRequest(client, "GET", payload, link, header, appendMode, isParameters)
	// Make POST Request
	makeRequest(client, "POST", payload, link, header, appendMode, isParameters)
	// Make OPTIONS Request
	makeRequest(client, "OPTIONS", payload, link, header, appendMode, isParameters)
	// Make PUT Request
	makeRequest(client, "PUT", payload, link, header, appendMode, isParameters)
}

func makeRequest(client *http.Client, method string, payload string, link string, header string, appendMode bool, isParameters bool) {
	fmt.Printf(NoticeColor, "\n[*] Making request with ", method)
	fmt.Println("")
	if isParameters {
		u, err := url.Parse(link)
		if err != nil {
			return
		}
		qs := url.Values{}
		for param, vv := range u.Query() {
			if appendMode {
				fmt.Printf(TextColor, "[*] Parameter:  ", param)
				qs.Set(param, vv[0]+payload)
			} else {
				fmt.Printf(TextColor, "[*] Parameter:  ", param)
				qs.Set(param, payload)
			}
		}
		u.RawQuery = qs.Encode()
		link = u.String()
	}
	fmt.Printf(InfoColor, "[-] Testing:  ", link)
	request, err := http.NewRequest(method, link, nil)
	if err != nil {
		return
	}
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) == 2 {
		request.Header.Set(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
	}
	client.Do(request)
}
