package main

import (
	"github.com/toshi0607/dfcx/internal/cli"
	"log"
	"os"
)

const (
	exitError = 1
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered, error: %v", err)
			os.Exit(exitError)
		}
	}()
	if err := cli.Run(os.Args); err != nil {
		os.Exit(exitError)
	}
}
