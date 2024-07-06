package main

import (
	"log"
	"os"

	"github.com/oxiginedev/hng11-stage-two-task/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatalf("something bad happened: %v", err)
	}

	os.Exit(0)
}
