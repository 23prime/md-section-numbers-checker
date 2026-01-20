package main

import (
	"flag"
	"fmt"

	"github.com/23prime/md-section-numbers-checker/internal/hello"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("md-section-numbers-checker %s\n", Version)
		fmt.Printf("  commit: %s\n", GitCommit)
		fmt.Printf("  built:  %s\n", BuildDate)
		return
	}

	greeting := hello.Greet("World")
	fmt.Println(greeting)
}
