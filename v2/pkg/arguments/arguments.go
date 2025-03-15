package arguments

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethicalhackingplayground/bxss/v2/pkg/colours"
)

type Arguments struct {
	Concurrency     int
	Header          string
	HeaderFile      string
	Payload         string
	PayloadFile     string
	Method          string
	AppendMode      bool
	Parameters      bool
	Debug           bool
	RateLimit       float64
	FollowRedirects bool
}

// Flag variables
var (
	debug           bool
	concurrency     int
	payload         string
	payloadFile     string
	method          string
	header          string
	headerFile      string
	appendMode      bool
	parameters      bool
	rateLimit       float64
	followRedirects bool
)

// ValidateArgs validates the arguments passed to the program and prints the
func (a *Arguments) ValidateArgs() {
	// prints the default help and exits the program with status 1.

	// The banner
	fmt.Printf(colours.BannerColor, `
  ____                     
 | __ )  __  __  ___   ___ 
 |  _ \  \ \/ / / __| / __|
 | |_) |  >  <  \__ \ \__ \
 |____/  /_/\_\ |___/ |___/
                                        
	`, "")
	fmt.Printf(colours.TextColor, "", "v0.0.3")
	fmt.Printf("\n")

	// Check if at least one header and one payload option is provided
	if (a.Header == "" && a.HeaderFile == "") && (a.Payload == "" && a.PayloadFile == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// NewArguments parses the command line flags and returns a pointer to an Arguments
// object. The Arguments object contains the parsed values of the flags, which can
// be used to configure the program. The function will panic if there is an error
// parsing the flags.
func NewArguments() *Arguments {

	// Define the flags
	flag.IntVar(&concurrency, "c", 30, "Set the concurrency level for the scanner")
	flag.StringVar(&header, "H", "", "Set a single custom header to test for blind XSS")
	flag.StringVar(&headerFile, "hf", "", "Path to file containing headers to test for blind XSS")
	flag.StringVar(&payload, "p", "", "The blind XSS payload to test")
	flag.StringVar(&payloadFile, "pf", "", "Path to file containing payloads to test for blind XSS")
	flag.BoolVar(&appendMode, "a", false, "Append the payload to the parameter value when testing")
	flag.BoolVar(&parameters, "t", false, "Test the parameters for blind XSS by appending the payload to the parameter value")
	flag.StringVar(&method, "X", "", "The HTTP method to use when testing (GET, POST, etc.)")
	flag.BoolVar(&debug, "v", false, "Enable debug mode to view full request details and debug information")
	flag.Float64Var(&rateLimit, "rl", 0, "Rate limit in requests per second (optional to prevent abuse)")
	flag.BoolVar(&followRedirects, "f", false, "Follow redirects when testing (optional)")

	// Parse the arguments
	flag.Parse()

	return &Arguments{
		Concurrency:     concurrency,
		Header:          header,
		HeaderFile:      headerFile,
		Payload:         payload,
		PayloadFile:     payloadFile,
		Method:          method,
		AppendMode:      appendMode,
		Parameters:      parameters,
		Debug:           debug,
		RateLimit:       rateLimit,
		FollowRedirects: followRedirects,
	}
}
