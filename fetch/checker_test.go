package fetch

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestMain(m *testing.M) {
	err := installPlaywright()
	if err != nil {
		slog.Error("failed to install playwright", slog.Any("error", err))
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestChecker_check(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectTitle string
		expectText  string
	}{
		{
			name:        "available",
			input:       "testdata/available",
			expectTitle: "Stock status",
			expectText:  "Stock is available",
		},
		{
			name:        "unavailable",
			input:       "testdata/unavailable",
			expectTitle: "",
			expectText:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			fs := http.FileServer(http.Dir(tt.input))
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				slog.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
				fs.ServeHTTP(w, r)
			}))

			var mu sync.Mutex
			gotTitle := ""
			gotText := ""

			c := NewChecker(CheckerConfig{
				URL:       srv.URL,
				Frequency: 5 * time.Millisecond,
				notifier: notifierFunc(func(title string, text string, iconPath string, urgency string) error {
					mu.Lock()
					defer mu.Unlock()
					gotTitle = title
					gotText = text

					cancel()
					return nil
				}),
				expect: playwright.NewPlaywrightAssertions(500),
			})

			err := c.Run(ctx)
			assert.Assert(t, err)

			t.Run("Checker", func(t *testing.T) {
				mu.Lock()
				defer mu.Unlock()

				assert.Check(t, cmp.Equal(tt.expectTitle, gotTitle))
				assert.Check(t, cmp.Equal(tt.expectText, gotText))
			})
		})
	}
}

type notifierFunc func(title string, text string, iconPath string, urgency string) error

func (n notifierFunc) Push(title string, text string, iconPath string, urgency string) error {
	return n(title, text, iconPath, urgency)
}
