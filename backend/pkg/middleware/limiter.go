package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit is middleware for handling rate limiter request per ip address
//
// first params is requestPerSecond
//
// second params is for burst,
//
// This burst is like a basket containing N tokens (N being the burst value). Each HTTP request retrieves a token from this basket. If the basket runs out of tokens, the rate limit occurs.
//
// The first parameter, RequestPerSecond, refills the basket with the specified number of tokens every second. If the basket reaches N, the token refill will stop.
func RateLimit(requestsPerSecond int, burst int) func(handler http.Handler) http.Handler {

	var mu sync.Mutex
	clients := make(map[string]*client)

	// set go routing to delete clients that have been inactive for a few minutes
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, c := range clients {
				// delete client if lastSeen > 2 minute
				if time.Since(c.lastSeen) > 2*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get client ip address
			ip := r.RemoteAddr

			mu.Lock()
			if _, found := clients[ip]; !found {
				// create new ratelimiter client if not exist
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
				}
			}
			// update lastSeen
			clients[ip].lastSeen = time.Now()

			// check if allow
			if !clients[ip].limiter.Allow() {
				// need to unlock if request not allowed
				mu.Unlock()
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// need to unlock if request allow
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
