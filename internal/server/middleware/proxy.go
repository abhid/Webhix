package middleware

import (
	"net"
	"net/http"
	"strings"
)

type TrustedProxies struct {
	ranges []*net.IPNet
}

func NewTrustedProxies(cidrs []string) *TrustedProxies {
	tp := &TrustedProxies{}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			return nil
		}

		tp.ranges = append(tp.ranges, ipNet)
	}

	return tp
}

func (tp *TrustedProxies) BehindProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tp.isTrusted(r) {
			http.Error(w, "untrusted proxy", http.StatusForbidden)
			return
		}

		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")
			clientIP := strings.TrimSpace(parts[0])

			if net.ParseIP(clientIP) != nil {
				r.RemoteAddr = net.JoinHostPort(clientIP, "0")
			}
		}

		if proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); proto != "" {
			r.URL.Scheme = proto
		}

		if host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); host != "" {
			r.Host = host
			r.URL.Host = host
		}

		next.ServeHTTP(w, r)
	})
}

func (tp *TrustedProxies) isTrusted(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	for _, ipNet := range tp.ranges {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
