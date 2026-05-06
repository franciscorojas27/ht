package main

import (
	"context"
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

	yml, ls, name := InitConfig()

	data, err := os.ReadFile(yml)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	requests, err := parseYML(data)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}

	if ls {
		if yml == "" {
			fmt.Println("Please provide a YAML file using the -yml flag.")
			return
		}
		requests.ListRequests()
		return
	}

	var req HTRequest
	if name != "" {
		found := false
		for _, r := range requests.Requests {
			if r.Name == name {
				req = r
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Request with name '%s' not found.\n", name)
			return
		}
	} else {
		if len(requests.Requests) == 0 {
			fmt.Println("No requests found in the YAML file.")
			return
		}
		req = requests.Requests[0]
	}

	if requests.Config.BaseURL != "" {
		req.URL = requests.Config.BaseURL + req.URL
	}

	var timeout time.Duration
	if requests.Config.Timeout > 0 {
		timeout = time.Duration(requests.Config.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqClient, err := http.NewRequestWithContext(ctx, req.Method, req.URL, nil)

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	LoadHeaders(reqClient, requests.Config.Headers)
	LoadHeaders(reqClient, req.Headers)

	client := &http.Client{}
	timeInit := time.Now()
	resp, err := client.Do(reqClient)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	dumpHeaders, err := httputil.DumpResponse(resp, false)
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
		return
	}

	color.Blue("HT [%s] %s", req.Method, req.URL)
	TimeTaken(timeInit)
	fmt.Printf("%s", dumpHeaders)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if strings.Contains(contentType, "application/json") {
		formattedJSON, err := FormatJSON(bodyBytes)
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(formattedJSON)
	} else {
		fmt.Printf("%s", bodyBytes)
	}
}
