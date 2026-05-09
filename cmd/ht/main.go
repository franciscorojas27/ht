package main

import (
	"fmt"
	"ht/internal/ht"
	"ht/internal/ui"
	"os"
)

func main() {
	flags := ui.InitFlags()
	if err := ht.Run(flags); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
