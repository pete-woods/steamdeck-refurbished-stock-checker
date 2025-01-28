package fetch

import (
	"github.com/playwright-community/playwright-go"
)

var (
	runOptions = &playwright.RunOptions{
		Browsers: []string{"chromium-headless-shell"},
	}
)

func installPlaywright() error {
	return playwright.Install(runOptions)
}

func startPlaywright(rawurl string) (playwright.Page, func() error, error) {
	var cleanups cleaner

	pw, err := playwright.Run(runOptions)
	if err != nil {
		return nil, cleanups.Cleanup, err
	}
	cleanups.Add(pw.Stop)

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return nil, cleanups.Cleanup, err
	}
	cleanups.Add(func() error {
		return browser.Close()
	})

	size := &playwright.Size{
		Width:  1920,
		Height: 1500,
	}
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Screen:   size,
		Viewport: size,
	})
	if err != nil {
		return nil, cleanups.Cleanup, err
	}
	cleanups.Add(func() error {
		return page.Close()
	})

	_, err = page.Goto(rawurl)
	if err != nil {
		return nil, cleanups.Cleanup, err
	}

	return page, cleanups.Cleanup, nil
}
