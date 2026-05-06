package ui

import (
	"fmt"
	"encoding/json"
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
