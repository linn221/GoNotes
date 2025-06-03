package main

import (
	"linn221/shop/config"
	"linn221/shop/models"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	DB                *gorm.DB
	Cache             *config.RedisCache
	ImageDirectory    string
	Readers           *models.ReadServices
	Port              string
	GeneralRateLimit  func(http.Handler) http.Handler
	ResourceRateLimit func(http.Handler) http.Handler
}
