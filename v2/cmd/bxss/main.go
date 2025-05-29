package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ethicalhackingplayground/bxss/v2/pkg/arguments"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/payloads"
	"golang.org/x/time/rate"
)

var (
	limiter *rate.Limiter
	args    *arguments.Arguments
)

func main() {

	// Create the arguments
	args = arguments.NewArguments()
	if args == nil {
		return
	}

	// Validate the arguments
	args.ValidateArgs()

	if args.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(args.RateLimit), 1)
	}
	// Create the payload parser
	payloadParser := payloads.NewPayload(args)
	if payloadParser == nil {
		fmt.Printf(colours.ErrorColor, "Error creating payload parser: "+"Something went wrong")
		os.Exit(1)
	}

	// Handle custom request file if specified
	if args.RequestFile != "" {
		fmt.Printf(colours.InfoColor, "Using custom request file: "+args.RequestFile)

		// Create request parser from the payloads package
		requestParser := payloads.NewRequestParser(args.RequestFile)
		if requestParser == nil {
			fmt.Printf(colours.ErrorColor, "Error creating request parser for file: "+args.RequestFile)
			os.Exit(1)
		}

		// Process the custom requests
		var payloadList []string
		if args.Payload != "" {
			payloadList = []string{args.Payload}
		}

		err := requestParser.ProcessCustomRequests(limiter, payloadList)
		if err != nil {
			fmt.Printf(colours.ErrorColor, "Error processing custom requests: "+err.Error())
			os.Exit(1)
		}

		// Exit after processing custom requests
		fmt.Printf(colours.SuccessColor, "Custom requests processed successfully.")
		os.Exit(0)
	}

	var headers []string
	if args.HeaderFile != "" {
		var err error
		headers, err = payloadParser.ReadLinesFromFile()
		if err != nil {
			fmt.Printf(colours.ErrorColor, "Error reading header file: "+err.Error())
			return
		}
	} else if args.Header != "" {
		headers = []string{args.Header}
	}

	var payloads []string
	if args.PayloadFile != "" {
		var err error
		payloads, err = payloadParser.ReadLinesFromFile()
		if err != nil {
			fmt.Printf(colours.ErrorColor, "Error reading payload file: "+err.Error())
			return
		}
	} else if args.Payload != "" {
		payloads = []string{args.Payload}
	}

	fmt.Printf(colours.NoticeColor, "Please Be Patient for bxss"+"")

	// Create a channel to send work items to the worker pool
	workChan := make(chan string)

	// Create a worker pool
	var wg sync.WaitGroup
	for i := 0; i < args.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for link := range workChan {
				payloadParser.ProcessPayloadsAndHeaders(limiter, link, payloads, headers)
			}
		}()
	}

	// Start sending the work items to the channel
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // Increase buffer size
		for scanner.Scan() {
			link := strings.TrimSpace(scanner.Text())
			if link == "" {
				continue // Skip empty lines
			}
			workChan <- link
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf(colours.ErrorColor, "Error reading input: "+err.Error())
		}
		close(workChan)
	}()

	wg.Wait()

	// Log completion message
	fmt.Printf(colours.SuccessColor, "Scan completed successfully.")
	fmt.Println("")
}
