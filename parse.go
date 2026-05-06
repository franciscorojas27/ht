package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func parseEnvYML(data []byte, vars HTVars) []byte {
	mapper := func(key string) string {
		if val, ok := vars[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return fmt.Sprintf("${%s}", key)
	}

	output := os.Expand(string(data), mapper)
	return []byte(output)
}

func parseYML(data []byte) (HT, error) {
	var tempVars struct {
		Vars HTVars `yaml:"vars"`
	}

	yaml.Unmarshal(data, &tempVars)

	if tempVars.Vars != nil {
		data = parseEnvYML(data, tempVars.Vars)
	}

	var requests HT
	if err := yaml.Unmarshal(data, &requests); err != nil {
		fmt.Printf("Error parsing final YAML: %v\n", err)
		return HT{}, err
	}
	return requests, nil
}
