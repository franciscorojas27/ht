package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"ht/internal/ht"
	"ht/internal/runner"
	"ht/internal/ui"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {

	yml, ls, name, dryRun, verbose := ui.InitFlags()

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

	bodyReader, _, err := targetReq.PrepareBody()
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

	if dryRun {
		ui.DryRun(reqClient)
		return
	}
	InitTime := time.Now()

	res, err := runner.DoRequest(reqClient, InitTime, verbose)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	finalTime := res.FinalTime

	dnsDur := res.DNSDur
	connDur := res.ConnDur
	tlsDur := res.TLSDur
	firstByteDur := res.FirstByte
	tlsVersion := res.TLSVersion
	sawTLS := res.SawTLS

	dumpHeaders := res.DumpHeaders

	color.Blue("HT [%s] %s", targetReq.Method, targetReq.URL)
	ui.TimeTaken(finalTime)
	if verbose {
		f := func(d time.Duration) string {
			if d == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
		}
		tlsVerStr := "-"
		switch tlsVersion {
		case tls.VersionTLS13:
			tlsVerStr = "v1.3"
		case tls.VersionTLS12:
			tlsVerStr = "v1.2"
		case tls.VersionTLS11:
			tlsVerStr = "v1.1"
		case tls.VersionTLS10:
			tlsVerStr = "v1.0"
		}
		fmt.Printf("◌ DNS Lookup      : %s\n", f(dnsDur))
		fmt.Printf("◌ TCP Connection  : %s\n", f(connDur))
		if sawTLS {
			fmt.Printf("◌ TLS Handshake   : %s (%s)\n", f(tlsDur), tlsVerStr)
		} else {
			fmt.Printf("◌ TLS Handshake   : -\n")
		}
		fmt.Printf("◌ First Byte      : %s\n", f(firstByteDur))
		fmt.Println(strings.Repeat("─", 35))
	}
	fmt.Printf("%s", dumpHeaders)
	fmt.Println(strings.Repeat("-", 35))
	bodyBytes := res.Body
	if strings.Contains(res.Headers.Get("Content-Type"), "application/json") {
		formattedJSON, err := ui.FormatJSON(bodyBytes)
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Printf("%s\n", formattedJSON)
	} else {
		fmt.Printf("%s", bodyBytes)
	}
}
