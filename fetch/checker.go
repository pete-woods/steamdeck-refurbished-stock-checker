package fetch

import (
	"context"
	"log/slog"
	"time"

	"github.com/0xAX/notificator"
)

type Checker struct {
	notify *notificator.Notificator
	cfg    CheckerConfig
}

type CheckerConfig struct {
	URL       string
	Frequency time.Duration
}

func NewChecker(cfg CheckerConfig) *Checker {
	notify := notificator.New(notificator.Options{
		AppName: "Steam Deck Checker",
	})

	return &Checker{
		cfg:    cfg,
		notify: notify,
	}
}

func (c *Checker) Run(ctx context.Context) error {
	err := installPlaywright()
	if err != nil {
		return err
	}

	t := time.NewTicker(c.cfg.Frequency)
	done := ctx.Done()
	for {
		select {
		case <-done:
			slog.Info("Checker stopped")
			return nil
		case <-t.C:
			err := c.check()
			if err != nil {
				return err
			}
		}
	}
}

func (c *Checker) check() (err error) {
	slog.Debug("Checking steam deck...")

	page, cleanup, err := startPlaywright(c.cfg.URL)
	defer closer(cleanup, &err)
	if err != nil {
		return err
	}

	addToCart := page.GetByText("add to cart").First()
	err = expect.Locator(addToCart).ToBeVisible()
	if err != nil {
		slog.Debug("No stock found")
		return nil
	}

	slog.Info("Stock available")
	return c.notify.Push("Stock status", "Stock is available", "", notificator.UR_CRITICAL)
}
