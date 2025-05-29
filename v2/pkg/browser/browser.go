package browser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
)

// BrowserType represents the type of browser to use
type BrowserType string

const (
	// Chrome browser
	Chrome BrowserType = "chrome"

	// Chromium browser
	Chromium BrowserType = "chromium"

	// Firefox browser
	Firefox BrowserType = "firefox"
)

// Browser represents a browser instance
type Browser struct {
	Type     BrowserType
	Path     string
	browsers []string
}

// NewBrowser creates a new browser instance
func NewBrowser(browserType string, customPath string) *Browser {
	bt := BrowserType(browserType)
	if bt != Chrome && bt != Chromium && bt != Firefox {
		fmt.Printf(colours.ErrorColor, "Unsupported browser type: "+browserType+". Using Chrome as default.")
		bt = Chrome
	}

	b := &Browser{
		Type: bt,
		Path: customPath,
	}

	// Initialize possible browser paths
	b.initBrowserPaths()

	return b
}

// initBrowserPaths initializes the possible browser paths based on the browser type
func (b *Browser) initBrowserPaths() {
	switch b.Type {
	case Chrome:
		b.browsers = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/local/bin/google-chrome",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
		}
	case Chromium:
		b.browsers = []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/local/bin/chromium",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"C:\\Program Files\\Chromium\\Application\\chrome.exe",
		}
	case Firefox:
		b.browsers = []string{
			"/usr/bin/firefox",
			"/usr/local/bin/firefox",
			"/Applications/Firefox.app/Contents/MacOS/firefox",
			"C:\\Program Files\\Mozilla Firefox\\firefox.exe",
			"C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe",
		}
	}

	// If custom path is provided, add it to the beginning of the list
	if b.Path != "" {
		b.browsers = append([]string{b.Path}, b.browsers...)
	}
}

// findBrowserPath finds the path to the browser executable
func (b *Browser) findBrowserPath() (string, error) {
	// If a custom path is provided and it exists, use it
	if b.Path != "" {
		if _, err := os.Stat(b.Path); err == nil {
			return b.Path, nil
		}
		fmt.Printf(colours.WarningColor, "Custom browser path not found: "+b.Path+". Trying default locations.")
	}

	// Check each possible path
	for _, path := range b.browsers {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// If we're here, we couldn't find the browser
	return "", errors.New("browser executable not found")
}

// CreateContext creates a new browser context
func (b *Browser) CreateContext(ctx context.Context) (context.Context, context.CancelFunc, error) {
	path, err := b.findBrowserPath()
	if err != nil {
		// Provide helpful error message with installation instructions
		b.printBrowserInstallationHelp()
		return nil, nil, fmt.Errorf("browser not found: %w", err)
	}

	// Different browser types require different approaches
	switch b.Type {
	case Chrome, Chromium:
		// Create Chrome/Chromium context
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(path),
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-extensions", true),
			chromedp.Flag("disable-setuid-sandbox", true),
			chromedp.Flag("disable-web-security", true),
			chromedp.WindowSize(1920, 1080),
		)

		allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
		// Don't defer cancel here - the allocator context must live as long as the browser context
		
		browserCtx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(format string, args ...interface{}) {
			// Suppress chromedp logs unless in debug mode
		}))

		// Add a timeout for browser operations
		browserCtx, timeoutCancel := context.WithTimeout(browserCtx, 10*time.Second)
		// Don't defer here either - we're returning this context

		// Ensure browser is started
		if err := chromedp.Run(browserCtx, chromedp.Navigate("about:blank")); err != nil {
			allocCancel()
			timeoutCancel()
			return nil, nil, fmt.Errorf("failed to start browser: %w", err)
		}

		// Return a cancel function that cleans up both contexts
		combinedCancel := func() {
			timeoutCancel()
			allocCancel()
		}
		return browserCtx, combinedCancel, nil

	case Firefox:
		// Currently, chromedp doesn't support Firefox directly
		// For Firefox, we need to use a different approach or library
		// This is a placeholder - in a real implementation, you might use another library for Firefox
		return nil, nil, errors.New("firefox support is currently experimental and not fully implemented")
	}

	return nil, nil, errors.New("unsupported browser type")
}

// printBrowserInstallationHelp prints helpful instructions for installing the required browser
func (b *Browser) printBrowserInstallationHelp() {
	fmt.Printf(colours.ErrorColor, "Browser not found: "+string(b.Type))
	fmt.Println()

	switch b.Type {
	case Chrome:
		fmt.Printf(colours.InfoColor, "To install Google Chrome, you can use:")
		fmt.Println("\nDebian/Ubuntu:")
		fmt.Println("  wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb")
		fmt.Println("  sudo apt install ./google-chrome-stable_current_amd64.deb")
		fmt.Println("  sudo cp /usr/bin/google-chrome-stable /usr/bin/google-chrome")
		fmt.Println("\nFedora:")
		fmt.Println("  sudo dnf install https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm")
		fmt.Println("\nOr specify a custom path with --browser-path flag")

	case Chromium:
		fmt.Printf(colours.InfoColor, "To install Chromium, you can use:")
		fmt.Println("\nDebian/Ubuntu:")
		fmt.Println("  sudo apt install chromium-browser")
		fmt.Println("\nFedora:")
		fmt.Println("  sudo dnf install chromium")
		fmt.Println("\nOr specify a custom path with --browser-path flag")

	case Firefox:
		fmt.Printf(colours.InfoColor, "To install Firefox, you can use:")
		fmt.Println("\nDebian/Ubuntu:")
		fmt.Println("  sudo apt install firefox")
		fmt.Println("\nFedora:")
		fmt.Println("  sudo dnf install firefox")
		fmt.Println("\nOr specify a custom path with --browser-path flag")
	}

	fmt.Println()
}

// BrowserPool represents a pool of browser contexts
type BrowserPool struct {
	browser        *Browser
	pool           chan context.Context
	cancelFuncs    []context.CancelFunc
	maxWorkers     int
	mu             sync.Mutex
	ctx            context.Context
	cancel         context.CancelFunc
	initializing   bool
	initialized    bool
	initErrCount   int
	initialization sync.WaitGroup
}

// NewBrowserPool creates a new browser pool
func NewBrowserPool(browser *Browser, maxWorkers int) *BrowserPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &BrowserPool{
		browser:      browser,
		pool:         make(chan context.Context, maxWorkers),
		cancelFuncs:  make([]context.CancelFunc, 0, maxWorkers),
		maxWorkers:   maxWorkers,
		ctx:          ctx,
		cancel:       cancel,
		initializing: false,
		initialized:  false,
	}

	return pool
}

// Initialize initializes the browser pool
func (p *BrowserPool) Initialize() error {
	p.mu.Lock()
	if p.initializing || p.initialized {
		p.mu.Unlock()
		return nil
	}
	p.initializing = true
	p.initialization.Add(1)
	p.mu.Unlock()

	defer p.initialization.Done()

	fmt.Printf(colours.InfoColor, fmt.Sprintf("Initializing browser pool with %d workers...\n", p.maxWorkers))

	for i := 0; i < p.maxWorkers; i++ {
		browserCtx, cancel, err := p.browser.CreateContext(p.ctx)
		if err != nil {
			p.mu.Lock()
			p.initErrCount++
			count := p.initErrCount
			p.mu.Unlock()

			// Only log the first error to avoid spam
			if count == 1 {
				fmt.Printf(colours.WarningColor, fmt.Sprintf("Error initializing browser worker: %v\n", err))
			}
			continue
		}

		p.mu.Lock()
		p.cancelFuncs = append(p.cancelFuncs, cancel)
		p.pool <- browserCtx
		p.mu.Unlock()

		fmt.Printf(colours.SuccessColor, fmt.Sprintf("Browser worker %d initialized\n", i+1))
	}

	p.mu.Lock()
	p.initializing = false

	// Check if we actually initialized any browsers
	if len(p.cancelFuncs) > 0 {
		p.initialized = true
		p.mu.Unlock()
		fmt.Printf(colours.SuccessColor, fmt.Sprintf("Browser pool initialized with %d workers\n", len(p.cancelFuncs)))
		return nil
	}

	p.mu.Unlock()
	return fmt.Errorf("failed to initialize any browser workers")
}

// GetContext gets a browser context from the pool
func (p *BrowserPool) GetContext() (context.Context, error) {
	if !p.initialized && !p.initializing {
		err := p.Initialize()
		if err != nil {
			return nil, err
		}
	}

	// If we're still initializing, wait for it to complete
	if p.initializing {
		p.initialization.Wait()
	}

	// If we failed to initialize, create a one-time context
	if !p.initialized {
		fmt.Printf(colours.WarningColor, "Using one-time browser context as pool initialization failed\n")
		ctx, cancel, err := p.browser.CreateContext(p.ctx)
		if err != nil {
			return nil, err
		}

		// Since this is one-time use, we'll clean it up when done
		go func() {
			<-ctx.Done()
			cancel()
		}()

		return ctx, nil
	}

	// Normal pool operation
	select {
	case ctx := <-p.pool:
		return ctx, nil
	case <-p.ctx.Done():
		return nil, errors.New("browser pool is closed")
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout waiting for browser context")
	}
}

// ReleaseContext returns a browser context to the pool
func (p *BrowserPool) ReleaseContext(ctx context.Context) {
	if !p.initialized {
		// For one-time contexts, we don't return to the pool
		return
	}

	select {
	case p.pool <- ctx:
		// Successfully returned to pool
	case <-p.ctx.Done():
		// Pool is closed, don't return
	case <-time.After(1 * time.Second):
		// If we can't return it to the pool in a reasonable time, discard it
		fmt.Printf(colours.WarningColor, "Timeout returning browser context to pool, discarding\n")
	}
}

// Close closes the browser pool and all browser instances
func (p *BrowserPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cancel()

	// Wait a moment for contexts to clean up
	time.Sleep(100 * time.Millisecond)

	// Cancel all browser contexts
	for _, cancel := range p.cancelFuncs {
		cancel()
	}
	p.cancelFuncs = nil
	p.initialized = false
}

// RequestParser represents a parser for custom HTTP requests
type RequestParser struct {
	FilePath string
}

// NewRequestParser creates a new request parser
func NewRequestParser(filePath string) *RequestParser {
	return &RequestParser{
		FilePath: filePath,
	}
}

// ParseRequests parses the custom HTTP requests from the file
// The file format should be a simple text file with one request per line
// Each line should be in the format: METHOD URL [HEADER:VALUE]...
// Example: GET https://example.com User-Agent:CustomAgent X-Custom:Value
func (p *RequestParser) ParseRequests() ([]*http.Request, error) {
	if p.FilePath == "" {
		return nil, errors.New("no request file path provided")
	}

	// Open the file
	file, err := os.Open(p.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open request file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	var requests []*http.Request
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line into a request
		req, err := p.parseRequestLine(line, lineNum)
		if err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading request file: %w", err)
	}

	if len(requests) == 0 {
		return nil, errors.New("no valid requests found in file")
	}

	return requests, nil
}

// parseRequestLine parses a single line into an http.Request
func (p *RequestParser) parseRequestLine(line string, lineNum int) (*http.Request, error) {
	// Split the line into parts
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil, fmt.Errorf("line %d: invalid request format, expected at least METHOD and URL", lineNum)
	}

	// Extract method and URL
	method := strings.ToUpper(parts[0])
	url := parts[1]

	// Validate method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"HEAD": true, "OPTIONS": true, "PATCH": true,
	}
	if !validMethods[method] {
		return nil, fmt.Errorf("line %d: invalid HTTP method '%s'", lineNum, method)
	}

	// Create the request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("line %d: failed to create request: %w", lineNum, err)
	}

	// Add headers if provided
	for i := 2; i < len(parts); i++ {
		headerPart := parts[i]
		// Headers should be in the format Header:Value
		if strings.Contains(headerPart, ":") {
			headerParts := strings.SplitN(headerPart, ":", 2)
			if len(headerParts) == 2 {
				headerName := strings.TrimSpace(headerParts[0])
				headerValue := strings.TrimSpace(headerParts[1])
				req.Header.Add(headerName, headerValue)
			} else {
				return nil, fmt.Errorf("line %d: invalid header format '%s'", lineNum, headerPart)
			}
		} else {
			return nil, fmt.Errorf("line %d: invalid header format '%s'", lineNum, headerPart)
		}
	}

	return req, nil
}

// ExecuteRequests executes all parsed requests and returns the responses
func (p *RequestParser) ExecuteRequests(ctx context.Context) ([]*http.Response, error) {
	requests, err := p.ParseRequests()
	if err != nil {
		return nil, err
	}

	var responses []*http.Response
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, req := range requests {
		// Use the provided context
		reqWithCtx := req.WithContext(ctx)

		// Execute the request
		resp, err := client.Do(reqWithCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request to %s: %w", req.URL.String(), err)
		}

		responses = append(responses, resp)
	}

	return responses, nil
}
