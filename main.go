package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pete-woods/steamdeck-refurbished-stock-checker/fetch"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("Error fetching store:", err)
		return
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	c := fetch.NewChecker()
	return c.Run(ctx)
}
