package main

import (
	"context"
	"fmt"
	"ht/internal/ht"
	"ht/internal/ui"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {

	yml, ls, name := ui.InitFlags()

	data, err := os.ReadFile(yml)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	requests, err := ht.ParseYML(data)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}

	if ls {
		requests.ListRequests()
		return
	}

	targetReq, found := requests.FindRequest(name)
	if !found {
		return
	}

	if requests.Config.BaseURL != "" {
		targetReq.URL = requests.Config.BaseURL + targetReq.URL
	}

	var timeout time.Duration
	if requests.Config.Timeout > 0 {
		timeout = time.Duration(requests.Config.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}
	
	bodyReader, contentType, err := targetReq.PrepareBody()
	if err != nil {
		fmt.Printf("Error preparing request body: %v\n", err)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqClient, err := http.NewRequestWithContext(ctx, targetReq.Method, targetReq.URL, bodyReader)

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	ht.LoadHeaders(reqClient, requests.Config.Headers)
	ht.LoadHeaders(reqClient, targetReq.Headers)

	client := &http.Client{}
	timeInit := time.Now()
	resp, err := client.Do(reqClient)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	dumpHeaders, err := httputil.DumpResponse(resp, false)
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
		return
	}

	color.Blue("HT [%s] %s", targetReq.Method, targetReq.URL)
	ui.TimeTaken(timeInit)
	fmt.Printf("%s", dumpHeaders)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if strings.Contains(contentType, "application/json") {
		formattedJSON, err := ui.FormatJSON(bodyBytes)
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(formattedJSON)
	} else {
		fmt.Printf("%s", bodyBytes)
	}
}
