package middlewares

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"linn221/shop/utils"
	"net/http"
	"strconv"
)

type SessionMiddleware struct {
	Cache services.CacheService
}

// set userId in context if authenticated
func (m *SessionMiddleware) Middleware(next http.Handler) http.Handler {

	respondInvalidSession := func(w http.ResponseWriter, r *http.Request) {
		utils.RemoveTokenCookies(w)
		http.Error(w, "invalid session", http.StatusUnauthorized)
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		// remove cookies to avoid infinite loop
		// http.SetCookie(w, &http.Cookie{
		// 	Name:    "token",
		// 	Expires: time.Unix(0, 0), // Set to past
		// 	MaxAge:  -1,              // Also ensures deletion
		// 	Path:    "/",
		// 	Domain:  "",
		// })
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookies, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token := cookies.Value
		//2d check token length
		val, err := m.Cache.GetH("Token:"+token, "userId")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId, err := strconv.Atoi(val)
		if err != nil || userId < 1 {
			respondInvalidSession(w, r)
			return
		}

		ctx := myctx.SetUserId(r.Context(), userId)
		ctx = myctx.SetAuth(ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
