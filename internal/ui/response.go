package ui

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/franciscorojas27/HT/internal/runner"

	"github.com/fatih/color"
)

func RenderResponse(method string, url string, res *runner.RequestResult, verbose bool) error {
	finalTime := res.FinalTime

	dnsDur := res.DNSDur
	connDur := res.ConnDur
	tlsDur := res.TLSDur
	firstByteDur := res.FirstByte
	tlsVersion := res.TLSVersion
	sawTLS := res.SawTLS

	dumpHeaders := res.DumpHeaders

	color.Blue("HT [%s] %s", method, url)
	TimeTaken(finalTime)
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
		formattedJSON, err := FormatJSON(bodyBytes)
		if err != nil {
			return fmt.Errorf("error formatting JSON: %w", err)
		}
		fmt.Printf("%s\n", formattedJSON)
	} else {
		fmt.Printf("%s", bodyBytes)
	}

	return nil
}
