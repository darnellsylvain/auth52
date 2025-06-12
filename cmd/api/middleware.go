package api

import (
	"log"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
	"golang.org/x/time/rate"
)

type rateLimiterClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func (api *API) rateLimiter(next http.Handler) http.Handler {
	go func() {
		for {
			time.Sleep(time.Minute)

			api.rateLimiterMu.Lock()
			for ip, client := range api.rateLimiterClients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(api.rateLimiterClients, ip)
				}
			}
			api.rateLimiterMu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if api.config.Limiter.Enabled {
			ip := GetIPAddressFromRequest(r).String()

			api.rateLimiterMu.Lock()

			if _, found := api.rateLimiterClients[ip]; !found {
				api.rateLimiterClients[ip] = &rateLimiterClient{
					limiter: rate.NewLimiter(rate.Limit(api.config.Limiter.RPS), api.config.Limiter.Burst),
				}
			}

			api.rateLimiterClients[ip].lastSeen = time.Now()

			if !api.rateLimiterClients[ip].limiter.Allow() {
				api.rateLimiterMu.Unlock()
				api.rateLimitExceededResponse(w, r)
				return
			}

			api.rateLimiterMu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) requireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			api.unauthorizedResponse(w, r, err)
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			api.unauthorizedResponse(w, r, err)
			return
		}

		ctx := auth.SetClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (api *API) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				w.Header().Set("Connection", "close")
				sendJSON(w, http.StatusInternalServerError, err, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
