package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	yml := flag.String("yml", "", "YAML file containing the request configuration")
	ls := flag.Bool("ls", false, "List all requests in the YAML file")
	name := flag.String("name", "", "Name of the request to execute (if not specified, the first request will be executed)")

	flag.Parse()

	data, err := os.ReadFile(*yml)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	requests, err := parseYML(data)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}

	if *ls {
		if *yml == "" {
			fmt.Println("Please provide a YAML file using the -yml flag.")
			return
		}
		requests.ListRequests()
		return
	}

	var req1 HTRequest
	if *name != "" {
		found := false
		for _, r := range requests.Requests {
			if r.Name == *name {
				req1 = r
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Request with name '%s' not found.\n", *name)
			return
		}
	} else {
		if len(requests.Requests) == 0 {
			fmt.Println("No requests found in the YAML file.")
			return
		}
		req1 = requests.Requests[0]
	}

	var bodyReader io.Reader
	if req1.Body != nil {
		switch v := req1.Body.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		default:
			b, merr := json.Marshal(v)
			if merr != nil {
				fmt.Printf("Error marshaling body: %v\n", merr)
				return
			}
			bodyReader = bytes.NewReader(b)
			if req1.Headers == nil {
				req1.Headers = make(map[string]string)
			}
			if _, ok := req1.Headers["Content-Type"]; !ok {
				req1.Headers["Content-Type"] = "application/json"
			}
		}
	}

	if requests.Config.BaseURL != "" {
		req1.URL = requests.Config.BaseURL + req1.URL
	}

	req, _ := http.NewRequest(req1.Method, req1.URL, bodyReader)

	LoadHeaders(req, requests.Config.Headers)
	LoadHeaders(req, req1.Headers)

	client := &http.Client{}
	timeInit := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	var dumpIn []byte
	if strings.Contains(contentType, "application/json") {
		dumpIn, err = httputil.DumpResponse(resp, false)
	} else {
		dumpIn, err = httputil.DumpResponse(resp, true)
	}
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
		return
	}
	color.Blue("HT [%s] %s", req1.Method, req1.URL)
	color.Yellow("Time taken: %s", time.Since(timeInit))
	fmt.Printf("%s\n", dumpIn)

	if strings.Contains(contentType, "application/json") {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		formattedJSON, err := FormatJSON(bodyBytes)
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(formattedJSON)
	}
}
