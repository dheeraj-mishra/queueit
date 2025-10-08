package middleware

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight (OPTIONS) requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Access-Control-Allow-Origin", "*") // allow all origins
	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 	if r.Method == "OPTIONS" {
	// 		w.WriteHeader(http.StatusOK)
	// 		return
	// 	}
	// 	next.ServeHTTP(w, r)
	// })
}
