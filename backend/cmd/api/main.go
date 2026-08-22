package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rafawastaken/tick-storm/backend/internal/app"
)

func main() {
	// Turns Ctrl+C and SIGTERM (what Docker and Kubernetes send) into a
	// cancelled context.
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	defer a.Close()

	if err := a.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
