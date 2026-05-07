package ui

import "flag"

func InitFlags() (string, bool, string, bool, bool) {
	yml := flag.String("yml", "", "YAML file configuration")
	ls := flag.Bool("ls", false, "List all requests")
	name := flag.String("name", "", "Request name to execute")
	dryRun := flag.Bool("dry-run", false, "Print the request without executing it")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()
	return *yml, *ls, *name, *dryRun, *verbose
}
