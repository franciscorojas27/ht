package runner

import (
	"crypto/tls"
	"net/http"
)

func NewClient(insecure bool) *http.Client {
	if !insecure {
		return &http.Client{}
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &http.Client{Transport: transport}
}
