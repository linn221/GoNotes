package middlewares

import (
	"errors"
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
	"strconv"
	"strings"
)

type SessionMiddleware struct {
	Cache services.CacheService
}

// set userId in context if authenticated
func (m *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Token")
		if token == "" {
			next.ServeHTTP(w, r)
			// http.Error(w, "token is required", http.StatusUnauthorized)
			return
		}

		//2d check token length
		val, ok, err := m.Cache.GetValue("Token:" + token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}
		userId, shopId, err := extractIdsFromCache(val)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := myctx.SetIds(r.Context(), userId, shopId)
		ctx = myctx.SetAuth(ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func extractIdsFromCache(s string) (int, string, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return 0, "", errors.New("error splitting the ids")
	}
	userId, err := strconv.Atoi(ss[0])
	if err != nil {
		return 0, "", err
	}
	shopId := ss[1]

	return userId, shopId, nil
}
