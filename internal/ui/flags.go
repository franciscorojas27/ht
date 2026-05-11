package ui

type Flags struct {
	YML      string
	LS       bool
	Name     string
	DryRun   bool
	Verbose  bool
	URL      string
	Method   string
	Headers  []string
	Data     string
	HeadOnly bool
	Insecure bool
	User     string
}
