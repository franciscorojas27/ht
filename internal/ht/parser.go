package ht

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseEnv(data []byte, vars map[string]any) []byte {
	mapper := func(key string) string {
		if val, ok := vars[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return fmt.Sprintf("${%s}", key)
	}

	output := os.Expand(string(data), mapper)
	return []byte(output)
}

func LoadEnvFile(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envs := make(map[string]any)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			envs[key] = val
		}
	}
	return envs, scanner.Err()
}

func ParseYML(data []byte) (HT, error) {
	var temp struct {
		Vars HTVars `yaml:"vars"`
		Env  string `yaml:"env_file"`
	}

	if err := yaml.Unmarshal(data, &temp); err != nil {
		return HT{}, fmt.Errorf("initial unmarshal error: %w", err)
	}

	vars := make(map[string]any)

	if temp.Env != "" {
		envs, err := LoadEnvFile(temp.Env)
		if err != nil {
			return HT{}, fmt.Errorf("error loading env_file: %w", err)
		}
		maps.Copy(vars, envs)
	}

	if temp.Vars != nil {
		maps.Copy(vars, temp.Vars)
	}

	if len(vars) > 0 {
		data = ParseEnv(data, vars)
	}

	var requests HT
	if err := yaml.Unmarshal(data, &requests); err != nil {
		return HT{}, fmt.Errorf("final unmarshal error: %w", err)
	}

	return requests, nil
}
