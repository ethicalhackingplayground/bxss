package payloads

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ethicalhackingplayground/bxss/v2/pkg/arguments"
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
	newScanner := scan.NewScanner(limiter, p.args.RateLimit, p.args.FollowRedirects, p.args.AppendMode, p.args.Parameters, p.args.Debug, p.args.ShowTimestamp)
	link = p.EnsureProtocol(link)
	fmt.Printf(colours.NoticeColor, "[+] Checking Url Scheme: ", link)
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
