package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	cfg := fetch.CheckerConfig{}

	flag.StringVar(&cfg.URL, "url", "https://store.steampowered.com/sale/steamdeckrefurbished/", "URL to check")
	flag.DurationVar(&cfg.Frequency, "frequency", 15*time.Minute, "time between checks")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	c := fetch.NewChecker(cfg)
	return c.Run(ctx)
}
