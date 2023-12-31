package am

import "net/http"

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: Improve this error handling
				http.Error(w, "Internal Server ErrorMsg", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
