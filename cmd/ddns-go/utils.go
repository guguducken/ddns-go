package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func Usage() func() {
	return func() {
		command := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "ddns-go is a library written in golang for managing DNS records of domain names and dynamic IP addresses.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [Options]\n\n", command)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -log-level=debug\n", command)
	}
}
