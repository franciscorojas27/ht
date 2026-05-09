package ui

import "flag"

type HeaderList []string

func (h *HeaderList) String() string {
	if h == nil {
		return ""
	}
	return ""
}

func (h *HeaderList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

type Flags struct {
	YML      string
	LS       bool
	Name     string
	DryRun   bool
	Verbose  bool
	URL      string
	Method   string
	Headers  HeaderList
	Data     string
	HeadOnly bool
	Insecure bool
	User     string
}

func InitFlags() Flags {
	var headers HeaderList

	yml := flag.String("yml", "", "YAML file configuration")
	ls := flag.Bool("ls", false, "List all requests")
	name := flag.String("name", "", "Request name to execute")
	dryRun := flag.Bool("dry-run", false, "Print the request without executing it")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	url := flag.String("url", "", "Quick URL to hit")
	method := flag.String("X", "", "HTTP method to use with the quick URL")
	flag.Var(&headers, "H", "Add a header to the request (can be repeated)")
	data := flag.String("d", "", "Request body data")
	headOnly := flag.Bool("I", false, "Fetch headers only (HEAD)")
	insecure := flag.Bool("k", false, "Allow insecure TLS")
	user := flag.String("u", "", "Basic auth user:pass")

	flag.Parse()
	return Flags{
		YML:      *yml,
		LS:       *ls,
		Name:     *name,
		DryRun:   *dryRun,
		Verbose:  *verbose,
		URL:      *url,
		Method:   *method,
		Headers:  headers,
		Data:     *data,
		HeadOnly: *headOnly,
		Insecure: *insecure,
		User:     *user,
	}
}
