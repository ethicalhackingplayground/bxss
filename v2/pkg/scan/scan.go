package scan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/browser"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
	"golang.org/x/time/rate"
)

type ScannerInterface interface {
	Scan(url string)
}

type ScannerConfig struct {
	AppendMode      bool
	IsParameters    bool
	RateLimit       float64
	Method          string
	FollowRedirects bool
	Limiter         *rate.Limiter
	Debug           bool
	Trace           bool
	BrowserType     string
	BrowserPath     string
	WorkerPool      int
	RequestFile     string
}

type Scanner struct {
	Config      ScannerConfig
	Client      *http.Client
	browserPool *browser.BrowserPool
	mu          sync.Mutex
}

func NewScanner(limiter *rate.Limiter, config *ScannerConfig) *Scanner {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !config.FollowRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Default to Chrome if no browser type specified
	browserType := config.BrowserType
	if browserType == "" {
		browserType = "chrome"
	}

	// Create browser instance
	b := browser.NewBrowser(browserType, config.BrowserPath)

	// Create a browser pool with the specified number of workers
	workerCount := config.WorkerPool
	if workerCount <= 0 {
		workerCount = 2 // Default to 2 workers if not specified
	}

	browserPool := browser.NewBrowserPool(b, workerCount)

	// Initialize the browser pool in the background
	go func() {
		err := browserPool.Initialize()
		if err != nil {
			fmt.Printf(colours.WarningColor, fmt.Sprintf("Warning: Browser pool initialization failed, will use one-time contexts: %v\n", err))
		}
	}()

	return &Scanner{
		Config:      *config,
		Client:      client,
		browserPool: browserPool,
	}
}

// Scan sends HTTP requests with different methods to a specified URL using a given payload and header.
// It respects the rate limiting if a limiter is set, pausing for a short duration between requests.
// The function iterates over a list of HTTP methods (GET, POST, OPTIONS, PUT) and invokes MakeRequest
// for each, using the provided payload and header. The function outputs the header and payload details
// to the console in a colored format.
func (s *Scanner) Scan(url string, payload string, header string) {
	fmt.Println("================================================================================")
	if s.Config.Limiter != nil {
		s.Config.Limiter.Wait(context.Background())
	}
	time.Sleep(500 * time.Microsecond)
	fmt.Println("")

	if header != "" {
		fmt.Printf(colours.InfoColor, "Using Header: "+header)
	}
	if s.Config.Trace {
		payload = strings.Replace(payload, "{LINK}", url, 1)
		fmt.Printf(colours.InfoColor, "**Using Trace Mode**"+"")
		fmt.Printf(colours.InfoColor, "New Payload:"+payload)
		fmt.Printf("\n")
	} else {
		fmt.Printf(colours.InfoColor, "Using Payload: "+payload)
		fmt.Printf("\n")
	}

	if s.Config.Method != "" {
		// Split the list of methods seperated with a comma if the comma exists
		// Otherwise just use the method passed
		if strings.Contains(s.Config.Method, ",") {
			methods := strings.Split(s.Config.Method, ",")
			for _, method := range methods {
				s.MakeRequest(method, payload, url, header, s.Config.AppendMode, s.Config.IsParameters)
			}
		} else {
			s.MakeRequest(s.Config.Method, payload, url, header, s.Config.AppendMode, s.Config.IsParameters)
		}
	} else {
		methods := []string{"GET", "POST", "OPTIONS", "PUT"}
		for _, method := range methods {
			s.MakeRequest(method, payload, url, header, s.Config.AppendMode, s.Config.IsParameters)
		}
	}

	fmt.Println("================================================================================")
}

// setheaders returns a task list that sets the passed headers.
func (s *Scanner) Setheaders(host string, headers map[string]interface{}, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate(host),
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
	fmt.Printf(colours.NoticeColor, "Method: "+method)

	u, err := url.Parse(link)
	if err != nil {
		fmt.Printf(colours.InfoColor, "Error parsing URL: "+err.Error())
		return
	}

	if isParameters {
		qs := u.Query()
		for param, vv := range qs {
			if appendMode {
				fmt.Printf(colours.NoticeColor, "Parameter: "+param)
				qs.Set(param, vv[0]+payload)
			} else {
				fmt.Printf(colours.NoticeColor, "Parameter: "+param)
				qs.Set(param, payload)
			}

		}
		u.RawQuery = qs.Encode()
	}

	fmt.Printf(colours.NoticeColor, ""+u.String()+"\n")
	request, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error creating request: "+err.Error())
		return
	}

	// Get a browser context from the pool instead of creating a new one each time
	ctx, err := s.getBrowserContext()
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error getting browser context: "+err.Error())
		return
	}
	defer s.releaseBrowserContext(ctx)

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
				// If appendMode is true, append the payload to the existing value
				if appendMode {
					request.Header.Set("User-Agent", headerValue+payload)
				} else {
					request.Header.Set("User-Agent", payload)
				}
			} else {
				// If appendMode is true, append the payload to the existing value
				if appendMode {
					request.Header.Set(headerName, headerValue+payload)
				} else {
					request.Header.Set(headerName, payload)
				}
			}
		} else {
			// If no value is provided, use the payload as the value
			request.Header.Set(header, payload)
		}

		// Get the headers from the request
		headers := make(map[string]interface{})
		for key := range request.Header {
			header := request.Header.Get(key)
			if s.Config.Debug {
				fmt.Printf(colours.DebugColor, "Header: "+key)
				fmt.Printf(colours.DebugColor, "Value: "+header)
			}

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
			fmt.Printf(colours.ErrorColor, "Error making request: "+err.Error())
			return
		}

	} else {
		err = chromedp.Run(ctx, chromedp.Navigate(u.String()))
		if err != nil {
			fmt.Printf(colours.ErrorColor, "Error making request: "+err.Error())
			return
		}
	}

	// Get the response from the request
	response, err := s.Client.Do(request)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error making request: "+err.Error())
		return
	}
	defer response.Body.Close()

	if s.Config.Debug {
		s.DebugRequest(request)
		s.DebugResponse(response)
	}

}

// DebugRequest dumps the request to the console in a human-readable format
// if s.Debug is true.
func (s *Scanner) DebugRequest(req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error dumping request: "+err.Error())
	} else {
		fmt.Printf("%s", "\n--- Request ---\n"+string(dump))
	}
}

// DebugResponse dumps the response to the console in a human-readable format
// if s.Debug is true. This includes the response headers and the response body.
func (s *Scanner) DebugResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error dumping response: "+err.Error())
	} else {
		fmt.Println(string(dump))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(colours.ErrorColor, "Error reading response body: "+err.Error())
	} else {
		fmt.Println(string(body))
	}
}

// getBrowserContext gets a browser context from the pool
func (s *Scanner) getBrowserContext() (context.Context, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.browserPool == nil {
		// Create a one-time context if no pool is available
		fmt.Printf(colours.WarningColor, "Browser pool not available, creating one-time context\n")
		ctx, _ := chromedp.NewContext(context.Background())
		return ctx, nil
	}

	// Get context from the pool
	ctx, err := s.browserPool.GetContext()
	if err != nil {
		// Fall back to creating a new context if the pool fails
		fmt.Printf(colours.WarningColor, fmt.Sprintf("Failed to get context from pool: %v, creating one-time context\n", err))
		oneTimeCtx, _ := chromedp.NewContext(context.Background())
		return oneTimeCtx, nil
	}

	return ctx, nil
}

// releaseBrowserContext returns a browser context to the pool
func (s *Scanner) releaseBrowserContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.browserPool != nil {
		s.browserPool.ReleaseContext(ctx)
	}
}

// Close cleans up resources used by the scanner
func (s *Scanner) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.browserPool != nil {
		s.browserPool.Close()
		s.browserPool = nil
	}
}
