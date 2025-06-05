package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guguducken/ddns-go/pkg/utils"
)

func Usage() func() {
	return func() {
		command := filepath.Base(os.Args[0])
		utils.MustWriteStringTo(os.Stdout, "ddns-go is a library written in golang for managing DNS records of domain names and dynamic IP addresses.\n\n")
		utils.MustWriteStringTo(os.Stdout, fmt.Sprintf("Usage: %s [Options]\n\n", command))
		utils.MustWriteStringTo(os.Stdout, "Options:\n")
		flag.PrintDefaults()
		utils.MustWriteStringTo(os.Stdout, "\nExamples:\n")
		utils.MustWriteStringTo(os.Stdout, fmt.Sprintf("  %s --config config.yaml\n", command))
	}
}

func PrintFigLet() {
	figlet := `       __    __
  ____/ /___/ /___  _____   ____ _____
 / __  / __  / __ \/ ___/  / __ \/ __ \
/ /_/ / /_/ / / / (__  )  / /_/ / /_/ /
\__,_/\__,_/_/ /_/____/   \__, \/____/
                         /____/

`
	utils.MustWriteStringTo(os.Stdout, figlet)
}
