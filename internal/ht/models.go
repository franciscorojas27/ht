package ht

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
	Name     string            `yaml:"name"`
	Method   string            `yaml:"method"`
	URL      string            `yaml:"url"`
	Headers  map[string]string `yaml:"headers"`
	Body     any               `yaml:"body"`
	BodyFile string            `yaml:"body_file"`
}

