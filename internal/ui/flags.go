package ui

import "flag"

func InitFlags() (string, bool, string) {
	yml := flag.String("yml", "", "YAML file configuration")
	ls := flag.Bool("ls", false, "List all requests")
	name := flag.String("name", "", "Request name to execute")

	flag.Parse()
	return *yml, *ls, *name
}
