package ht

import (
	"net/http"
)

func LoadHeaders(req *http.Request, headers map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}
