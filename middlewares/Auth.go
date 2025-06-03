package middlewares

import (
	"linn221/shop/myctx"
	"net/http"
)

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !myctx.IsAuth(r.Context()) {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
