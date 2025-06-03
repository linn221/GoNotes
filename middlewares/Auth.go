package middlewares

import (
	"linn221/shop/myctx"
	"net/http"
)

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !myctx.IsAuth(r.Context()) {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, r)
	})
}
