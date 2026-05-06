package main

import (
	"github.com/fatih/color"
)

type HT struct {
	Requests []HTRequest `yaml:"requests"`
	Config   HTConfig    `yaml:"config"`
}

type HTConfig struct {
	BaseURL string            `yaml:"base_url"`
	Timeout int               `yaml:"timeout"`
	Headers map[string]string `yaml:"headers"`
}

type HTVars map[string]any

type HTRequest struct {
	Name    string            `yaml:"name"`
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Body    any               `yaml:"body"`
}

func (rs HT) ListRequests() {
	for i, req := range rs.Requests {
		color.Cyan("Request #%d:", i+1)
		color.Yellow("Method: %s", req.Method)
		color.Yellow("URL: %s", req.URL)
	}
}
