package middlewares

import (
	"linn221/shop/services"
	"net/http"

	"github.com/go-chi/cors"
)

func Default(h http.Handler, cache services.CacheService) http.Handler {
	sessionMd := SessionMiddleware{Cache: cache}
	corsMiddleware := cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins:   []string{"*"}, // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	return Recovery(corsMiddleware(sessionMd.Middleware(Logging(h))))
}
