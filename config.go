package main

import (
	"flag"
)

func InitConfig() (string, bool, string) {
	yml := flag.String("yml", "", "YAML file containing the request configuration")
	ls := flag.Bool("ls", false, "List all requests in the YAML file")
	name := flag.String("name", "", "Name of the request to execute (if not specified, the first request will be executed)")

	flag.Parse()
	return *yml, *ls, *name
}