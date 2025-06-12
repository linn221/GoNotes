package main

import (
	"context"
	"linn221/shop/config"
	"linn221/shop/models"
	"linn221/shop/views"
	"path/filepath"
)

func main() {
	ctx := context.Background()
	db := config.ConnectDB()
	redisCache := config.ConnectRedis(ctx)
	dir := config.GetImageDirectory()
	port := config.GetPortNo()
	services := models.NewServices(db, redisCache)
	templates := views.NewEngine(config.GetBaseDir())
	assetDir := filepath.Join(config.GetBaseDir(), "../../files")

	// // rate limiting crud endpoints by userId
	// resourceRateLimit := middlewares.NewRateLimiter(redisCache.GetClient(), time.Minute*5, 2000, "r", func(r *http.Request) (string, error) {
	// 	ctx := r.Context()
	// 	userId, _, err := myctx.GetIdsFromContext(ctx)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return fmt.Sprint(userId), nil
	// })
	// // rate limit by IP address for all routes
	// generalRateLimit := middlewares.NewRateLimiter(redisCache.GetClient(), time.Minute*2, 300, "g", func(r *http.Request) (string, error) {
	// 	ip := utils.GetClientIP(r)
	// 	return ip, nil
	// })
	// serves http server
	app := &App{
		DB:             db,
		Cache:          redisCache,
		ImageDirectory: dir,
		Services:       services,
		Templates:      templates,
		Port:           port,
		AssetDirectory: assetDir,
		// ResourceRateLimit: resourceRateLimit,
		// GeneralRateLimit:  generalRateLimit,
	}
	app.Serve()
}
