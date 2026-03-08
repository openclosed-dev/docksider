package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"

	"github.com/openclosed-dev/docksider/internal/cmd"
)

//go:embed version.txt
var version string

func main() {

	log.SetFlags(0)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := cmd.NewRootCmd(version).ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}
