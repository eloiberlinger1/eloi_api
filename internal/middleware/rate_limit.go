package middleware

import (
	"net/http"
	"sync"
	"time"
)

type client struct {
	lastSeen time.Time
	count    int
}

var clients = make(map[string]*client)
var mu sync.Mutex

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()

		c, exists := clients[ip]
		if !exists {
			clients[ip] = &client{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(c.lastSeen) > 10*time.Second {
			delete(clients, ip)
			mu.Unlock()

			next.ServeHTTP(w, r)
			return
		}

		if c.count >= 5 {
			mu.Unlock()
			http.Error(w, "Trop de requêtes. Veuillez réessayer plus tard.", http.StatusTooManyRequests)
			return
		}

		c.count++
		c.lastSeen = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
