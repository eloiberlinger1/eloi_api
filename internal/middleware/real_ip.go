package middleware

import (
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(clientIP)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RemoteAddr = getClientIP(r)
		next.ServeHTTP(w, r)
	})
}
