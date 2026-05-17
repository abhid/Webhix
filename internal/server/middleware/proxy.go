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
			if clientIP := tp.forwardedClientIP(forwardedFor); clientIP != nil {
				r.RemoteAddr = net.JoinHostPort(clientIP.String(), "0")
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

func (tp *TrustedProxies) forwardedClientIP(forwardedFor string) net.IP {
	parts := strings.Split(forwardedFor, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		ip := net.ParseIP(strings.TrimSpace(parts[i]))
		if ip != nil && !tp.isTrustedIP(ip) {
			return ip
		}
	}

	return nil
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

	return tp.isTrustedIP(ip)
}

func (tp *TrustedProxies) isTrustedIP(ip net.IP) bool {
	for _, ipNet := range tp.ranges {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
