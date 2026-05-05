package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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
	var req HTRequest
	err = yaml.Unmarshal(data, &req)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}
	fmt.Printf("Parsed Request:\nMethod: %s\nURL: %s\nHeaders: %v\nBody: %s\n", req.Method, req.URL, req.Headers, req.Body)

}
