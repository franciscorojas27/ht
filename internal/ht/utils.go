package ht

import (
	"fmt"
	"github.com/fatih/color"
)

func (h HT) FindRequest(name string) (HTRequest, bool) {
	if name != "" {
		for _, r := range h.Requests {
			if r.Name == name {
				return r, true
			}
		}
		fmt.Printf("Request with name '%s' not found.\n", name)
		return HTRequest{}, false
	}

	if len(h.Requests) == 0 {
		fmt.Println("No requests found in the YAML file.")
		return HTRequest{}, false
	}

	return h.Requests[0], true
}

func (rs HT) ListRequests() {
	for i, req := range rs.Requests {
		color.Cyan("Request #%d:", i+1)
		color.Yellow(" Method: %s", req.Method)
		color.Yellow(" URL: %s", req.URL)
	}
}
