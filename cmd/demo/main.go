package main

import (
	"fmt"
	"linn221/shop/app"
	"linn221/shop/models"
	"linn221/shop/views"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("demo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Image{},
		&models.Label{},
		&models.Note{})
	if err != nil {
		panic("Error migrating: " + err.Error())
	}

	fmt.Println("connecting sqlite done")
	return db
}

func main() {

	// get the base directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// ctx := context.Background()
	db := connectSqlite()
	cache := NewInMemoryCache()
	imageDir := filepath.Join(dir, "images")
	port := "8080"
	services := models.NewServices(db, cache)
	templates := views.NewTemplates(filepath.Join(dir, "../../views/templates"))
	assetDir := filepath.Join(dir, "../../files/static")

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
	app := &app.App{
		DB:             db,
		Cache:          cache,
		ImageDirectory: imageDir,
		Services:       services,
		Templates:      templates,
		Port:           port,
		AssetDirectory: assetDir,
		// ResourceRateLimit: resourceRateLimit,
		// GeneralRateLimit:  generalRateLimit,
	}
	app.Serve()
}
