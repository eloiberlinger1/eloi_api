package middleware

import (
	"net"
	"net/http"
	"strings"
)

// getClientIP extracts the real client IP using X-Forwarded-For
func getClientIP(r *http.Request) string {
	// 1. On vérifie d'abord si le proxy nous a transmis l'IP
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For peut contenir plusieurs IP si la requête a traversé plusieurs proxys
		// La première IP de la liste est celle du client d'origine
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// 2. Fallback si on attaque l'API en direct sans proxy
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// RealIP is a custom middleware that replaces the r.RemoteAddr with the
// real client IP extracted by getClientIP.
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RemoteAddr = getClientIP(r)
		next.ServeHTTP(w, r)
	})
}
