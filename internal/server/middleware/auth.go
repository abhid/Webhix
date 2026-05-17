package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

type Auth struct {
	password  string
	secretKey string
}

func NewAuth(password, secretKey string) *Auth {
	return &Auth{
		password:  password,
		secretKey: secretKey,
	}
}

func (a *Auth) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/r/") {
			next.ServeHTTP(w, r)
			return
		}

		if a.isAuthorized(r) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="webhix"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

func (a *Auth) isAuthorized(r *http.Request) bool {
	if a.secretKey != "" {
		if key := r.Header.Get("X-Webhix-Key"); key != "" {
			if subtle.ConstantTimeCompare([]byte(key), []byte(a.secretKey)) == 1 {
				return true
			}
		}
		if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			key := strings.TrimPrefix(auth, "Bearer ")
			if subtle.ConstantTimeCompare([]byte(key), []byte(a.secretKey)) == 1 {
				return true
			}
		}
	}

	if a.password != "" {
		_, pass, ok := r.BasicAuth()
		if ok && subtle.ConstantTimeCompare([]byte(pass), []byte(a.password)) == 1 {
			return true
		}
	}

	return false
}
