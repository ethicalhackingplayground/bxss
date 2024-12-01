package scan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
	"golang.org/x/time/rate"
)

type ScannerInterface interface {
	Scan(url string)
}

type Scanner struct {
	AppendMode      bool
	IsParameters    bool
	RateLimit       float64
	FollowRedirects bool
	Client          *http.Client
	Limiter         *rate.Limiter
	Debug           bool
	ShowTimestamp   bool
}

func NewScanner(limiter *rate.Limiter, rateLimit float64, followRedirects bool, appendMode, parameters bool, debug bool, showTimeStamp bool) *Scanner {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !followRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return &Scanner{
		AppendMode:      appendMode,
		IsParameters:    parameters,
		RateLimit:       rateLimit,
		FollowRedirects: followRedirects,
		Client:          client,
		Limiter:         limiter,
		Debug:           debug,
		ShowTimestamp:   showTimeStamp,
	}
}

// Scan sends HTTP requests with different methods to a specified URL using a given payload and header.
// It respects the rate limiting if a limiter is set, pausing for a short duration between requests.
// The function iterates over a list of HTTP methods (GET, POST, OPTIONS, PUT) and invokes MakeRequest
// for each, using the provided payload and header. The function outputs the header and payload details
// to the console in a colored format.
func (s *Scanner) Scan(url string, payload string, header string) {
	fmt.Printf(colours.NoticeColor, "[+] Scanning ", url)
	fmt.Println("")

	if s.Limiter != nil {
		s.Limiter.Wait(context.Background())
	}
	time.Sleep(500 * time.Microsecond)
	fmt.Println("")
	if header != "" {
		fmt.Printf(colours.NoticeColor, "[+] Header:  ", header)
	}
	if payload != "" {
		fmt.Printf(colours.TextColor, "[+] Payload: ", payload)
	}
	fmt.Println("")

	methods := []string{"GET", "POST", "OPTIONS", "PUT"}
	for _, method := range methods {
		s.MakeRequest(method, payload, url, header, s.AppendMode, s.IsParameters)
	}
}

// setheaders returns a task list that sets the passed headers.
func (s *Scanner) Setheaders(host string, headers map[string]interface{}, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate(host),
		chromedp.Text(`#result`, res, chromedp.ByID, chromedp.NodeVisible),
	}
}

// MakeRequest constructs and sends an HTTP request with the specified method,
// payload, and headers to the given link. If isParameters is true, the payload
// is appended to each query parameter. The function also allows setting custom
// headers and handles special cases for the User-Agent header. It uses chromedp
// to modify the request and navigate to the link. If ShowTimestamp is true, a
// timestamp is printed. If Debug is true, the request and response are dumped
// to the console. The function returns no value.
func (s *Scanner) MakeRequest(method string, payload string, link string, header string, appendMode, isParameters bool) {
	fmt.Printf(colours.NoticeColor, "\n[*] Making request with ", method)
	fmt.Println("")

	u, err := url.Parse(link)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "\nError parsing URL: ", err.Error())
		return
	}

	if isParameters {
		qs := u.Query()
		for param, vv := range qs {
			if appendMode {
				fmt.Printf(colours.TextColor, "[*] Parameter:  ", param)
				qs.Set(param, vv[0]+payload)
			} else {
				fmt.Printf(colours.TextColor, "[*] Parameter:  ", param)
				qs.Set(param, payload)
			}
		}
		u.RawQuery = qs.Encode()
	}

	fmt.Printf(colours.InfoColor, "[-] Testing:  ", u.String())
	request, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "\nError creating request: ", err.Error())
		return
	}

	// Use chromedp to navigate to the link and modify the request
	// based on the payload and header
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Check if the header is empty
	if header != "" {

		// Remove existing headers that we're testing
		request.Header.Del("User-Agent")
		request.Header.Del("X-Forwarded-Host")
		request.Header.Del("X-Forwarded-For")

		// Set the header with the payload
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) == 2 {
			headerName := strings.TrimSpace(headerParts[0])
			headerValue := strings.TrimSpace(headerParts[1])

			// Special handling for User-Agent header
			if strings.ToLower(headerName) == "user-agent" {
				request.Header.Set("User-Agent", headerValue+payload)
			} else {
				request.Header.Set(headerName, headerValue+payload)
			}
		} else {
			// If no value is provided, use the payload as the value
			request.Header.Set(header, payload)
		}

		if s.ShowTimestamp {
			fmt.Printf(colours.InfoColor, "\n[*] Timestamp: ", time.Now().Format(time.RFC3339))
		}

		if s.Debug {
			s.DebugRequest(request)
		}

		// Get the headers from the request
		headers := make(map[string]interface{})
		for key := range request.Header {
			header := request.Header.Get(key)
			fmt.Printf(colours.DebugColor, "-> Header: \n", key)
			fmt.Printf(colours.DebugColor, "-> Value: \n", header)
			headers[key] = header
		}

		// Set the headers for the request using chromedp
		var res string
		err = chromedp.Run(ctx, s.Setheaders(
			u.String(),
			headers,
			&res,
		))
		if err != nil {
			fmt.Printf(colours.ErrorColor, "\nError making request: ", err.Error())
			return
		}

	} else {
		err = chromedp.Run(ctx, chromedp.Navigate(u.String()))
		if err != nil {
			fmt.Printf(colours.ErrorColor, "\nError making request: ", err.Error())
			return
		}
	}

	// Get the response from the request
	response, err := s.Client.Do(request)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "\nError making request: ", err.Error())
		return
	}
	defer response.Body.Close()

	if s.Debug {
		s.DebugResponse(response)
	}

}

// DebugRequest dumps the request to the console in a human-readable format
// if s.Debug is true.
func (s *Scanner) DebugRequest(req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "\nError dumping request: ", err.Error())
	} else {
		fmt.Printf(colours.DebugColor, "\n--- Request ---\n", string(dump))
	}
}

// DebugResponse dumps the response to the console in a human-readable format
// if s.Debug is true. This includes the response headers and the response body.
func (s *Scanner) DebugResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "\nError dumping response: ", err.Error())
	} else {
		fmt.Printf(colours.DebugColor, "\n--- Response ---\n", string(dump))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error reading response body: ", err.Error())
	} else {
		fmt.Printf(colours.DebugColor, "\n--- Response Body ---\n", string(body))
	}
}
