package ht

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func (req HTRequest) PrepareBody() (io.Reader, string, error) {
	if req.BodyFile != "" {
		content, err := os.ReadFile(req.BodyFile)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(content), "", nil
	}

	if req.Body != nil {
		switch v := req.Body.(type) {
		case string:
			return strings.NewReader(v), "", nil
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, "", err
			}
			return bytes.NewReader(b), "application/json", nil
		}
	}

	return nil, "", nil
}
