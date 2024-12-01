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
	AppendMode      bool
	Parameters      bool
	Debug           bool
	RateLimit       float64
	FollowRedirects bool
	ShowTimestamp   bool
}

// Flag variables
var (
	debug           bool
	showTimestamp   bool
	concurrency     int
	payload         string
	payloadFile     string
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
	 |  _ \              
 	 | |_) |_  _____ ___ 
	 |  * <\ \/ / *_/ __|
	 | |_) |>  <\__ \__ \
	 |____//_/\_\___/___/
	                     
                    
	`, "-- Coded by @z0idsec -- \n")

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
	flag.IntVar(&concurrency, "concurrency", 30, "Set the concurrency")
	flag.StringVar(&header, "header", "", "Set a single custom header")
	flag.StringVar(&headerFile, "headerFile", "", "Path to file containing headers to test")
	flag.StringVar(&payload, "payload", "", "The blind XSS payload")
	flag.StringVar(&payloadFile, "payloadFile", "", "Path to file containing payloads to test")
	flag.BoolVar(&appendMode, "appendMode", false, "Append the payload to the parameter")
	flag.BoolVar(&parameters, "parameters", false, "Test the parameters for blind xss")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode to view full request details")
	flag.Float64Var(&rateLimit, "rl", 0, "Rate limit in requests per second (optional)")
	flag.BoolVar(&followRedirects, "r", false, "Follow redirects (optional)")
	flag.BoolVar(&showTimestamp, "ts", false, "Show timestamp for each request (optional)")

	// Parse the arguments
	flag.Parse()

	return &Arguments{
		Concurrency:     concurrency,
		Header:          header,
		HeaderFile:      headerFile,
		Payload:         payload,
		PayloadFile:     payloadFile,
		AppendMode:      appendMode,
		Parameters:      parameters,
		Debug:           debug,
		RateLimit:       rateLimit,
		FollowRedirects: followRedirects,
		ShowTimestamp:   showTimestamp,
	}
}
