package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darnellsylvain/auth52/internal/config"
)

func TestRateLimiter_AllowsAndBlocks(t *testing.T) {

	api := &API{
		config: &config.Config{
			Limiter: config.LimiterConfig{
				Enabled: true,
				RPS:     2, // 2 requests per second
				Burst:   2,
			},
		},
		rateLimiterClients: make(map[string]*rateLimiterClient),
	}

	// Dummy handler that returns 200 OK
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := api.rateLimiter(next)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:12345" // IP used as limiter key

	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}
	}

	// Third request should be rate-limited
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", rec.Code)
	}

	// Wait for limiter to refill
	time.Sleep(1 * time.Second)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK after cooldown, got %d", rec.Code)
	}
}
