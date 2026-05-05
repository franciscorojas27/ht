package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"gopkg.in/yaml.v3"
)

type HTRequest struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

func main() {

	data, err := os.ReadFile("dev.yml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}
	var req1 HTRequest
	err = yaml.Unmarshal(data, &req1)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}

	req, _ := http.NewRequest(req1.Method, req1.URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HT/1.0)")
	dumpOut, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
		return
	}
	fmt.Printf("Dumped Request:\n%s\n", dumpOut)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	dumpIn, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
		return
	}
	fmt.Printf("Dumped Response:\n%s\n", dumpIn)
}
