package fetch

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/playwright-community/playwright-go"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestChecker_check(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedStock bool
	}{
		{
			name:          "available",
			input:         "testdata/available",
			expectedStock: true,
		},
		{
			name:          "unavailable",
			input:         "testdata/unavailable",
			expectedStock: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := http.FileServer(http.Dir(tt.input))
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				slog.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
				fs.ServeHTTP(w, r)
			}))

			var notified atomic.Bool
			c := NewChecker(CheckerConfig{
				URL: srv.URL,
				notifier: notifierFunc(func(title string, text string, iconPath string, urgency string) error {
					notified.Store(true)
					return nil
				}),
				expect: playwright.NewPlaywrightAssertions(500),
			})

			err := c.check()
			assert.Assert(t, err)
			assert.Check(t, cmp.Equal(tt.expectedStock, notified.Load()))
		})
	}
}

type notifierFunc func(title string, text string, iconPath string, urgency string) error

func (n notifierFunc) Push(title string, text string, iconPath string, urgency string) error {
	return n(title, text, iconPath, urgency)
}
