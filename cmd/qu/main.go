package main

import (
	"context"
	"log"
	"os"

	"github.com/queeck/cli/internal/services/app"
)

func main() {
	program, err := app.Default()
	if err != nil {
		log.Printf("failed to make default app:\n%v\n", err)
		os.Exit(1)
	}
	if err = program.Run(context.Background()); err != nil {
		log.Printf("failed to run app:\n%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
