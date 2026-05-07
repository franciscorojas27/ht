package ui

import (
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func TimeTaken(finalTime time.Duration) {
	if finalTime > 2*time.Second {
		color.Red("Time taken: %s", finalTime)
	} else if finalTime > 1*time.Second {
		color.Yellow("Time taken: %s", finalTime)
	} else {
		color.Green("Time taken: %s", finalTime)
	}
}
func DryRun(req *http.Request) {
	dumpReq, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
		return
	}
	sDump := string(dumpReq)

	parts := strings.SplitN(sDump, "\r\n\r\n", 2)
	headers := parts[0]
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	color.Green("Dry run mode enabled. Request details:")
	fmt.Printf("%s Headers:\n", color.BlueString(">"))
	fmt.Println(headers)

	if body != "" {
		fmt.Printf("%s Body:\n%s\n", color.BlueString(">"), body)
	} else {
		fmt.Printf("%s No body.\n", color.BlueString(">"))
	}
	color.Yellow("Request not sent.")
}
