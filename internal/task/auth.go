package task

import "net/http"

// APIKeyMiddleware protects mutating endpoints when API_KEY is configured.
func APIKeyMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key != "" && r.Method != http.MethodGet && r.Header.Get("X-API-Key") != key {
				writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
