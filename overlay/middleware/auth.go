package middleware

import "net/http"

func Auth(token string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenParam := r.URL.Query().Get("token")
			if tokenParam != token {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
