package middleware

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PATCH, PUT, DELETE, OPTIONS")

		// ⭐ Cloud Run で OPTIONS を安定させるための最小追加
		if r.Method == http.MethodOptions {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(http.StatusNoContent) // ← 204 の方が Cloud Run が安定
			return
		}

		next.ServeHTTP(w, r)
	})
}
