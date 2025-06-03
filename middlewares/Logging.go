package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		h.ServeHTTP(w, r)
		fmt.Printf("%s %s %s\n", r.Method, r.URL, time.Since(now))
	})
}
