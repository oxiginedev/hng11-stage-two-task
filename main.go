package main

import (
	"log"

	"github.com/oxiginedev/hng11-stage-two-task/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatalf("something bad happened: %v", err)
	}
}
