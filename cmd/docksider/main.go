package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/openclosed-dev/docksider/internal/cmd"
)

func main() {

	log.SetFlags(0)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := cmd.NewRootCmd().ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}
