package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TylerBrock/colorjson"
)

func FormatJSON(data []byte) (string, error) {
	var obj any
	err := json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON response: %v\n", err)
		return "", err
	}
	f := colorjson.NewFormatter()
	f.Indent = 4
	s, _ := f.Marshal(obj)
	return string(s), nil
}

func LoadHeaders(req *http.Request, headers map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}
