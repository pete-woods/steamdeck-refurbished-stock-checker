package fetch

import (
	"context"
	"log/slog"
	"time"

	"github.com/0xAX/notificator"
	"github.com/playwright-community/playwright-go"
)

type Checker struct {
	cfg CheckerConfig
}

type notifier interface {
	Push(title string, text string, iconPath string, urgency string) error
}

type CheckerConfig struct {
	URL       string
	Frequency time.Duration

	// used in tests
	notifier notifier
	expect   playwright.PlaywrightAssertions
}

func NewChecker(cfg CheckerConfig) *Checker {
	if cfg.expect == nil {
		cfg.expect = playwright.NewPlaywrightAssertions()
	}
	if cfg.notifier == nil {
		cfg.notifier = notificator.New(notificator.Options{
			AppName: "Steam Deck Checker",
		})
	}
	return &Checker{
		cfg: cfg,
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
	err = c.cfg.expect.Locator(addToCart).ToBeVisible()
	if err != nil {
		slog.Debug("No stock found")
		return nil
	}

	slog.Info("Stock available")
	return c.cfg.notifier.Push("Stock status", "Stock is available", "", notificator.UR_CRITICAL)
}
