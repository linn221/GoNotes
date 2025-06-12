package main

import (
	"linn221/shop/config"
	"linn221/shop/models"
	"linn221/shop/views"

	"gorm.io/gorm"
)

type App struct {
	DB             *gorm.DB
	Cache          *config.RedisCache
	ImageDirectory string
	AssetDirectory string
	Services       *models.CrudServices
	Port           string
	Templates      *views.Templates
	// GeneralRateLimit  func(http.Handler) http.Handler
	// ResourceRateLimit func(http.Handler) http.Handler
}
