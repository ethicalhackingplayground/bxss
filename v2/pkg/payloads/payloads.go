package payloads

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ethicalhackingplayground/bxss/v2/pkg/arguments"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/browser"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/scan"
	"golang.org/x/time/rate"
)

type PayloadParser struct {
	args *arguments.Arguments
}

func NewPayload(args *arguments.Arguments) *PayloadParser {
	return &PayloadParser{
		args: args,
	}
}

// readLinesFromFile reads a file line by line and returns the lines as a slice of strings.
//
// The lines are trimmed of whitespace. If there is an error reading the file,
// that error is returned. Otherwise, the function returns a slice of strings
// and a nil error.
func (p *PayloadParser) ReadLinesFromFile() ([]string, error) {
	file, err := os.Open(p.args.PayloadFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// processPayloadsAndHeaders reads lines from standard input, and for each line,
// sends a request with each payload to each header to the specified link.
//
// The function takes in a slice of payloads, a slice of headers, and booleans
// indicating whether to append the payload to the parameter, whether to test
// parameters, and whether to follow redirects. It also takes a client object
// with a timeout and a redirect policy.
//
// If there is an error reading the input, that error is printed to standard
// error. Otherwise, the function prints nothing and returns no value.
func (p *PayloadParser) ProcessPayloadsAndHeaders(limiter *rate.Limiter, link string, payloads []string, headers []string) {
	config := &scan.ScannerConfig{
		AppendMode:      p.args.AppendMode,
		IsParameters:    p.args.Parameters,
		RateLimit:       p.args.RateLimit,
		Method:          p.args.Method,
		FollowRedirects: p.args.FollowRedirects,
		Debug:           p.args.Debug,
		Trace:           p.args.Trace,
		BrowserType:     p.args.BrowserType,
		BrowserPath:     p.args.BrowserPath,
		WorkerPool:      p.args.WorkerPool,
		RequestFile:     p.args.RequestFile,
	}
	newScanner := scan.NewScanner(limiter, config)
	link = p.EnsureProtocol(link)
	fmt.Printf(colours.NoticeColor, "Checking URL Scheme: "+link)
	fmt.Println("")
	if len(headers) == 0 {
		for _, payload := range payloads {
			newScanner.Scan(link, payload, "")
		}
	} else {
		for _, payload := range payloads {
			for _, header := range headers {
				newScanner.Scan(link, payload, header)
			}
		}
	}

}

// EnsureProtocol verifies that the provided link has a protocol prefix.
// If the link does not start with "http://" or "https://", it prepends "https://" to the link.
// The function trims any leading or trailing whitespace from the link before checking the protocol.
// It returns the modified or unmodified link with the appropriate protocol.
func (p *PayloadParser) EnsureProtocol(link string) string {
	link = strings.TrimSpace(link)
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		return "https://" + link
	}
	return link
}

// RequestParser is a wrapper around browser.RequestParser for handling custom requests
type RequestParser struct {
	args     *arguments.Arguments
	filePath string
}

// NewRequestParser creates a new request parser for custom requests
func NewRequestParser(filePath string) *RequestParser {
	return &RequestParser{
		args:     nil, // Not needed for direct file processing
		filePath: filePath,
	}
}

// ProcessCustomRequests processes custom requests from a file
func (p *RequestParser) ProcessCustomRequests(limiter *rate.Limiter, payloads []string) error {
	// Create the browser request parser
	parser := browser.NewRequestParser(p.filePath)

	// Create browser context for executing requests
	browserType := "chrome" // Default fallback
	if p.args != nil && p.args.BrowserType != "" {
		browserType = p.args.BrowserType
	}
	b := browser.NewBrowser(browserType, "")
	ctx, cancel, err := b.CreateContext(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create browser context: %w", err)
	}
	defer cancel()

	// Execute the requests
	fmt.Printf(colours.InfoColor, "Processing custom requests from file...")
	responses, err := parser.ExecuteRequests(ctx)
	if err != nil {
		return err
	}

	// Report on the responses
	fmt.Printf(colours.InfoColor, fmt.Sprintf("Processed %d custom requests successfully", len(responses)))

	// Clean up responses
	for _, resp := range responses {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	return nil
}
